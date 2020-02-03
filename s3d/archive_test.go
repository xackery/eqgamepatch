package s3d

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
// archives/snd1.pfs
portc_lp.wav
dbrdg_lp.wav
// archives/snd2.pfs
torch_lp.wav
doorwd_c.wav
doorwd_o.wav
jumpland.wav
rainloop.wav
spelcast.wav
spelhit1.wav
spelhit2.wav
spelhit3.wav
spelhit4.wav
spell_1.wav
spell_2.wav
spell_3.wav
spell_4.wav
spell_5.wav
stepcrch.wav
steprun2.wav
steprun3.wav
steprunl.wav
stepwlk2.wav
stepwlk3.wav
stepwlkl.wav
stepwlkr.wav
thunder1.wav
thunder2.wav
torch_lp.wav
wattrd_1.wav
wattrd_2.wav
wind_lp1.wav
watundlp.wav
// archives/snd3.pfs
// archives/snd4.pfs
// archives/snd5.pfs
drake_lp.wav
fairy_lp.wav
ghost_lp.wav
mimic_lp.wav
vampk_lp.wav
runlp.wav
djin_idle.wav
gorilp.wav
der_lp.wav
wngflplp.wav
drix_lp.wav
vampkidle.wav
snakelp.wav
imp_lp.wav
waterlp.wav
walk_lp.wav
// archives/snd6.pfs
fire_lp.wav
rumblelp.wav
spinnrlp.wav
steamlp.wav
lava_lp.wav
lava2_lp.wav
waterlp.wav
// archives/snd7.pfs
dbrdg_lp.wav
portc_lp.wav
// archives/snd8.pfs
arena_bg_lp.wav
crystal_big_lp.wav
fire_large_lp.wav
fire_small_lp.wav
flag_flap_lp.wav
forcefield_lp.wav
ghosts_lp.wav
heartbeat_lp.wav
jungle_lp.wav
kegbar_lp.wav
lake_ripple_lp.wav
ocean_waves_lp.wav
siren_song_lp.wav
water_brook_lp.wav
waterfall_big_lp.wav
waterfall_med_lp.wav
whisper_lp.wav
wind_cave_lp.wav
wind_corr_lp.wav
wind_lite_lp.wav
wind_strong_lp.wav
// archives/snd9.pfs
bear_idle.wav
bear_idle2.wav
bronto_idle.wav
centaur_idle.wav
coldain_mal_idle.wav
coldfem_idle.wav
coldfem_idle2.wav
dervish_idle.wav
dragon_robo_idle.wav
drakf_wool_idle.wav
drakm_wool_idle.wav
encharmor_idle.wav
faun_idle.wav
giant_frost_idle.wav
giant_storm_idle.wav
hag_idle.wav
hellcat_idle1.wav
hgriff_idle.wav
icefreti_idle.wav
manti_idle.wav
manti_idle2.wav
monkey_idle.wav
oposs_idle.wav
otter_idle.wav
owl_idle1.wav
owl_idle2.wav
panther_blk_idle.wav
rabbit_idle.wav
ratman_wht_idle.wav
rockman_idle.wav
spectre_cold_idle.wav
totem_idle.wav
turtle_huge_idle.wav
turtle_med_idle.wav
walrus_idle.wav
wolf_dire_idle.wav
wyvern_idle.wav
yak_idle.wav
// archives/snd10.pfs
gmnidle.wav
felidle.wav
unbidle.wav
fugidle.wav
seridle.wav
shmidle.wav
spridle.wav
netidle.wav
akmidle.wav
muhidle.wav
vacidle.wav
sgridle.wav
gmfidle.wav
tigidle.wav
sowidle.wav
zelidle.wav
welidle.wav
rnbidle.wav
rhpidle.wav
vpmidle.wav
schidle.wav
skeidle.wav
wetidle.wav
sknidle.wav
eelidle.wav
khaidle.wav
kesidle.wav
snnidle.wav
hsmidle.wav
lcridle.wav
aelidle.wav
aknidle.wav
sdmidle.wav
volidle.wav
srvidle.wav
thoidle.wav
shridle.wav
tegidle.wav
owbidle.wav
galidle.wav
// archives/snd11.pfs
combat_lp.wav
wind_desert_lp.wav
sword_forcfield_lp.wav
tavern_lp.wav
vampyres_evil_lp.wav
nightime_background02_lp.wav
flag_flap_lp.wav
water_underwater_lp.wav
wheel_pottery_lp.wav
wind_caverns_lp.wav
waves01_lp.wav
drums_deep_lp.wav
grimling_chant_lp.wav
flag_flap_lite_lp.wav
wind_soft_lp.wav
door_forcfield_lp.wav
crowd_townhall_lp.wav
fire_torch02_lp.wav
chanting_drums_lp.wav
snake_chant_lp.wav
bar_lp.wav
fairies_twinkle_lp.wav
fire_torch01_lp.wav
chanting_lp.wav
drums_lp.wav
worm_trap_lp.wav
fire_bonfire_lp.wav
wind_water_birds_lp.wav
river_big_lp.wav
arena_lp.wav
nightime_background_lp.wav
electric_arcs_lp.wav
saw_wood_lp.wav
slime_bg_lp.wav
waterfall_big_lp.wav
marching_lp.wav
mining_cart_lp.wav
tegi_jam_lp.wav
tent_flap_lp.wav
coloseum_crowd_lp.wav
scientist_lab_lp.wav
mosquitos_lp.wav
whispering_lp.wav
// archives/snd12.pfs
// archives/snd13.pfs
// archives/snd14.pfs
water_sway_lp.wav
mephit_fire_chant_lp.wav
waterfall_med_lp.wav
big_battle_lp.wav
part_click_lp.wav
thunder_distant_lp.wav
rain_wood_roof_lp.wav
hum_low_lp.wav
slargh_chant_lp.wav
undead_chant_lp.wav
forest_amb_lp.wav
rat_chant_lp.wav
sand_shift_lp.wav
manastorm_lp.wav
gear_grind_lp.wav
wheel_mtl_creak_lp.wav
wind_moan_lp.wav
rain_hvy_leaves_lp.wav
wind_chimes_lp.wav
clock_tick_lp.wav
banner_flap_lp.wav
sand_storm_lp.wav
scream_help.wav
river_pus_lp.wav
wind_trees_lp.wav
gear_whir_lp.wav
blacksmith_lp.wav
*/

func TestArchiveLoad(t *testing.T) {
	assert := assert.New(t)
	path := "archives/snd14.pfs"
	a, err := New(path)
	if !assert.NoError(err) {
		t.Fatal(err)
	}

	assert.NotNil(a)
	t.Log(a.Count(), "files found in", path)
	err = a.ExtractAll("out")
	if !assert.NoError(err) {
		t.Fatal(err)
	}

	err = a.Save("test.s3d")
	if !assert.NoError(err) {
		t.Fatal(err)
	}
}

func TestSoundFind(t *testing.T) {
	path := ""
	for i := 1; i < 17; i++ {
		path = fmt.Sprintf("archives/snd%d.pfs", i)
		New(path)
	}
}
