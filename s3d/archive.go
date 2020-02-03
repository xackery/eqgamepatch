package s3d

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/lunixbochs/struc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	fileNameListCRC = 0x61580AC9
	chunkSize       = 8192
)

var crcs []uint32

// Archive represents a s3d archive
type Archive struct {
	path        string
	mutex       sync.RWMutex
	fileEntries fileEntries
}

type fileEntry struct {
	name       string
	fileHeader fileHeader
}

func (fe *fileEntry) crc() uint32 {
	if fe.name == "" {
		return 0
	}

	crc := uint32(0)
	for _, c := range fe.name {
		crc = crcCalc(c, crc)
	}
	crc = crcCalc('\x00', crc)

	return crc
}

func crcCalc(c rune, crc uint32) uint32 {
	index := ((crc >> 24) ^ uint32(c)) & 0xFF
	crc = ((crc << 8) ^ crcs[index])
	return crc
}

type pFSHeader struct {
	Offset        uint32
	MagicCookie   [4]byte
	VersionNumber uint32
}

type directoryHeader struct {
	Count uint32
}

type byOffset []fileHeader

func (a byOffset) Len() int           { return len(a) }
func (a byOffset) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byOffset) Less(i, j int) bool { return a[i].Offset < a[j].Offset }

type fileEntries []fileEntry

type byCRC []fileEntry

func (a byCRC) Len() int           { return len(a) }
func (a byCRC) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCRC) Less(i, j int) bool { return a[i].crc() < a[j].crc() }

type dataBlock struct {
	CompressedLength uint32
	InflatedLength   uint32
}

type fileNameCount struct {
	Count uint32
}

type fileNameLength struct {
	Length uint32
}

func init() {

	// build crcs, based on https://github.com/Shendare/EQZip/blob/master/EQArchive.cs#L51
	for i := 0; i < 256; i++ {
		crc := uint32(i << 24)

		for round := 0; round < 8; round++ {
			if (crc & 0x80000000) != 0 {
				crc = ((crc << 1) ^ 0x04C11DB7)
			} else {
				crc = (crc << 1)
			}
		}

		crcs = append(crcs, crc)
	}
}

// New creates a new archive
func New(path string) (*Archive, error) {
	a := new(Archive)
	a.path = path

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	headerBytes := make([]byte, 12)
	header := pFSHeader{}
	_, err = f.Read(headerBytes)
	if err != nil {
		return nil, errors.Wrap(err, "read pfsheader")
	}
	buffer := bytes.NewBuffer(headerBytes)
	err = binary.Read(buffer, binary.LittleEndian, &header)
	if err != nil {
		return nil, errors.Wrap(err, "binary read pfsheader")
	}
	// fmt.Printf("Directory Header Offset: %X\n", header.Offset)

	// Validate header
	if string(header.MagicCookie[:]) != "PFS " {
		return nil, fmt.Errorf("pfs magic cookie")
	}
	if header.VersionNumber != 131072 {
		return nil, fmt.Errorf("unknown version number header != 131072")
	}

	directoryHeaderBytes := make([]byte, 4)
	directoryHeader := directoryHeader{}
	_, err = f.ReadAt(directoryHeaderBytes, int64(header.Offset))
	if err != nil {
		return nil, errors.Wrap(err, "directory header")
	}
	buffer = bytes.NewBuffer(directoryHeaderBytes)
	err = binary.Read(buffer, binary.LittleEndian, &directoryHeader)
	if err != nil {
		return nil, errors.Wrap(err, "binary directory header")
	}
	// fmt.Printf("File Count: %d\n", directoryHeader.Count)

	// Get file crcs, offsets, and checksums
	fileHeaders := make([]fileHeader, 0)
	fileNameHeader := fileHeader{}
	for i := 0; i < int(directoryHeader.Count); i++ {
		fileHeaderBytes := make([]byte, 12)
		fileHeader := fileHeader{}
		_, err = f.ReadAt(fileHeaderBytes, int64(int(header.Offset)+4+i*12))
		if err != nil {
			return nil, errors.Wrapf(err, "offset %d file header read", i)
		}
		buffer = bytes.NewBuffer(fileHeaderBytes)
		err = binary.Read(buffer, binary.LittleEndian, &fileHeader)
		if err != nil {
			return nil, errors.Wrapf(err, "offset %d byte file header read", i)
		}
		// fmt.Printf("Parsed data: %+v\n", directoryHeader)

		if fileHeader.CRC == uint32(fileNameListCRC) {
			fileNameHeader = fileHeader
			// fmt.Printf("Directory Header Found. CRC: %X Offset: %X Size: %X\n", fileHeader.CRC, fileHeader.Offset, fileHeader.Size)
		} else {
			fileHeaders = append(fileHeaders, fileHeader)
			// fmt.Printf("File Header Found. CRC: %X Offset: %X Size: %X\n", fileHeader.CRC, fileHeader.Offset, fileHeader.Size)
		}
	}

	// Sort the offsets
	sort.Sort(byOffset(fileHeaders))
	fmt.Println("//", path)

	// Get file names
	fileNamedataBlockBytes := make([]byte, 8)
	fileNamedataBlock := dataBlock{}
	_, err = f.ReadAt(fileNamedataBlockBytes, int64(fileNameHeader.Offset))
	if err != nil {
		return nil, errors.Wrap(err, "data block")
	}
	buffer = bytes.NewBuffer(fileNamedataBlockBytes)
	err = binary.Read(buffer, binary.LittleEndian, &fileNamedataBlock)
	if err != nil {
		return nil, errors.Wrap(err, "binary data block")
	}
	// TODO: Read more file name blocks if necessary
	// fmt.Printf("Filename Header Block Found. Compressed Length: %X Inflated Length: %X\n", fileNamedataBlock.CompressedLength, fileNamedataBlock.InflatedLenth)

	fileNameBytes := make([]byte, fileNamedataBlock.CompressedLength)
	_, err = f.ReadAt(fileNameBytes, int64(fileNameHeader.Offset+8))
	if err != nil {
		return nil, errors.Wrap(err, "file name")
	}
	buffer = bytes.NewBuffer(fileNameBytes)
	r, err := zlib.NewReader(buffer)
	if err != nil {
		return nil, errors.Wrap(err, "zlib reader")
	}
	defer r.Close()

	// File count
	fileNameCountBytes := make([]byte, 4)
	fileNameCount := fileNameCount{}
	_, err = r.Read(fileNameCountBytes)
	if err != nil {
		return nil, errors.Wrap(err, "file name count bytes")
	}
	buffer = bytes.NewBuffer(fileNameCountBytes)
	err = binary.Read(buffer, binary.LittleEndian, &fileNameCount)
	if err != nil {
		return nil, errors.Wrap(err, "binary file name count bytes")
	}
	// fmt.Printf("Found %d Files\n", fileNameCount.Count)

	for i := 0; i < int(fileNameCount.Count); i++ {
		// File length
		fileNameLengthBytes := make([]byte, 4)
		fileNameLength := fileNameLength{}
		_, err = r.Read(fileNameLengthBytes)
		if err != nil {

			if strings.Contains(path, "gequip.s3d") {
				log.Debug().Msgf("gequip offset %d skipped, known invalid offset", i)
				continue
			}
			return nil, errors.Wrap(err, "file name length")
		}
		buffer = bytes.NewBuffer(fileNameLengthBytes)
		err = binary.Read(buffer, binary.LittleEndian, &fileNameLength)
		if err != nil {
			return nil, errors.Wrap(err, "binary file name length")
		}
		// fmt.Printf("Filename %d Length Found: %X bytes\n", i+1, fileNameLength.Length)

		// File name
		fileNameEntryBytes := make([]byte, fileNameLength.Length)
		_, err = r.Read(fileNameEntryBytes)
		if err != nil && err != io.EOF {
			return nil, errors.Wrap(err, "read file name entry bytes")
		}
		fileName := string(bytes.Trim(fileNameEntryBytes, "\x00"))

		a.fileEntries = append(a.fileEntries, fileEntry{
			name:       fileName,
			fileHeader: fileHeaders[i],
		})
	}
	for _, e := range a.fileEntries {
		if strings.Contains(e.name, ".wav") && (strings.Contains(e.name, "lp.") || strings.Contains(e.name, "idle")) {
			fmt.Println(e.name)
		}
	}
	return a, nil
}

// Count returns file count
func (a *Archive) Count() int {
	a.mutex.RLock()
	count := len(a.fileEntries)
	a.mutex.RUnlock()
	return count
}

// Names returns every file
func (a *Archive) Names() []string {
	files := []string{}
	a.mutex.RLock()
	for _, f := range a.fileEntries {
		files = append(files, f.name)
	}
	a.mutex.RUnlock()
	return files
}

// ExtractAll will extract every file found within the loaded s3d to outPath
func (a *Archive) ExtractAll(outPath string) error {
	if err := os.MkdirAll(outPath, 0700); err != nil {
		return err
	}
	names := a.Names()
	for _, name := range names {
		if err := a.Extract(name, outPath); err != nil {
			return errors.Wrap(err, name)
		}
	}
	return nil
}

// Extract extracts fileName from loaded s3d to outPath
func (a *Archive) Extract(fileName string, outPath string) error {
	a.mutex.RLock()
	var entry fileEntry
	isFound := false
	for _, entry = range a.fileEntries {
		if strings.ToLower(entry.name) == strings.ToLower(fileName) {
			isFound = true
			break
		}
	}
	a.mutex.RUnlock()
	if !isFound {
		return fmt.Errorf("%s not found", fileName)
	}

	fa, err := os.Open(a.path)
	if err != nil {
		return err
	}
	defer fa.Close()

	// Extract file
	var inflated uint32
	f, err := os.Create(fmt.Sprintf("%s/%s", outPath, entry.name))
	if err != nil {
		return errors.Wrap(err, "create")
	}
	defer f.Close()

	// Read multiple blocks (this only reads a single 8k block and writes it out)
	_, err = fa.Seek(int64(entry.fileHeader.Offset), 0)
	if err != nil {
		return errors.Wrap(err, "seek")
	}

	buf := new(bytes.Buffer)
	for inflated < entry.fileHeader.InflatedSize {
		filedataBlockBytes := make([]byte, 8)
		filedataBlock := dataBlock{}
		_, err = fa.Read(filedataBlockBytes)
		if err != nil {
			return errors.Wrap(err, "read data block bytes")
		}
		buf = bytes.NewBuffer(filedataBlockBytes)
		err = binary.Read(buf, binary.LittleEndian, &filedataBlock)
		if err != nil {
			return errors.Wrap(err, "parse data block bytes")
		}
		// fmt.Printf("Data Block Compressed: %X Inflated: %X\n", filedataBlock.CompressedLength, filedataBlock.InflatedLenth)
		fileBytes := make([]byte, filedataBlock.CompressedLength)
		_, err = fa.Read(fileBytes)
		if err != nil {
			return errors.Wrap(err, "read bytes")
		}
		buf = bytes.NewBuffer(fileBytes)
		r2, err := zlib.NewReader(buf)
		if err != nil {
			return errors.Wrap(err, "read zlib")
		}
		if _, err = io.Copy(f, r2); err != nil {
			return errors.Wrap(err, "copy")
		}
		if err = r2.Close(); err != nil {
			return errors.Wrap(err, "close")
		}

		inflated += filedataBlock.InflatedLength
	}

	return nil
}

// Save writes the archive to provided outPath
func (a *Archive) Save(outPath string) error {
	/*if err := os.MkdirAll(outPath, 0700); err != nil {
		return err
	}*/

	w, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer w.Close()

	// step 1: order entries by CRC
	sort.Sort(byCRC(a.fileEntries))

	//TODO: step 2: directory of filenames

	//for uint32 i := 0; i <
	// step 3: build the file header
	//https://github.com/Shendare/EQZip/blob/master/EQArchive.cs

	header := pFSHeader{
		Offset:        12,                              //start with 12, to be determined
		MagicCookie:   [4]byte{0x50, 0x46, 0x53, 0x20}, //80, 70, 83, 32}, //PFS
		VersionNumber: 131072,
	}

	fa, err := os.Open(a.path)
	if err != nil {
		return err
	}
	defer fa.Close()

	for _, e := range a.fileEntries {
		if err = e.fileHeader.calculateSize(fa); err != nil {
			return errors.Wrapf(err, "calculate size of %s", e.name)
		}
		header.Offset += 4 + 4 + e.fileHeader.CompressedSize
	}

	// step 4: write the header file
	if err := struc.PackWithOrder(w, &header, binary.LittleEndian); err != nil {
		return errors.Wrap(err, "pack")
	}

	// step 5: compressed file chunks
	return nil
}
