package s3d

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type fileHeader struct {
	CRC            uint32
	Offset         uint32
	InflatedSize   uint32
	CompressedSize uint32
}

// calculate the sizes of inflated and compressed information
func (fh *fileHeader) calculateSize(fa *os.File) error {
	if _, err := fa.Seek(int64(fh.Offset), 0); err != nil {
		return errors.Wrap(err, "seek")
	}

	var inflated uint32
	for inflated < fh.InflatedSize {
		filedataBlockBytes := make([]byte, 8)
		filedataBlock := dataBlock{}
		if _, err := fa.Read(filedataBlockBytes); err != nil {
			return errors.Wrap(err, "read data block bytes")
		}
		buf := bytes.NewBuffer(filedataBlockBytes)
		if err := binary.Read(buf, binary.LittleEndian, &filedataBlock); err != nil {
			return errors.Wrap(err, "parse data block bytes")
		}
		// fmt.Printf("Data Block Compressed: %X Inflated: %X\n", filedataBlock.CompressedLength, filedataBlock.InflatedLenth)
		fileBytes := make([]byte, filedataBlock.CompressedLength)
		if _, err := fa.Read(fileBytes); err != nil {
			return errors.Wrap(err, "read bytes")
		}
		buf = bytes.NewBuffer(fileBytes)
		r2, err := zlib.NewReader(buf)
		if err != nil {
			return errors.Wrap(err, "read zlib")
		}
		if _, err := io.Copy(ioutil.Discard, r2); err != nil {
			return errors.Wrap(err, "copy")
		}
		if err := r2.Close(); err != nil {
			return errors.Wrap(err, "close")
		}

		inflated += filedataBlock.InflatedLength
	}

	return nil
}
