package priyadi

import (
	"regexp"
	"strings"
)

func Encode(s string) string {
	if len(s) == 0 {
		return ""
	}
	s = strings.ToLower(s)
	s = DiftongNormalize(s, `[^a-z]`, ``)
	/* 1. Ubah ejaan lama menjadi ejaan baru: ubah oe menjadi u, tj
	   menjadi c, dj menjadi j. Untuk menghindari false positive,
	   jangan ubah j menjadi y kecuali jika ada pengubahan ejaan
	   lama menjadi ejaan baru yang lain. Hati-hati dengan urutan
	   pengubahan, jangan sampai dj berubah menjadi y. */
	s = strings.ReplaceAll(s, "oe", "u")
	s = strings.ReplaceAll(s, "tj", "c")
	s = strings.ReplaceAll(s, "dj", "j")

	/* 2. Ganti konsonan yang berderet menjadi satu konsonan saja.
	   Misalnya ‘anni’ menjadi ‘ani’. */
	s = RemoveConsecutiveConsonants(s)

	/* 3. Normalkan diftong: ubah ai di akhir kata menjadi ay, au
	   di akhir kata menjadi aw dan oi di akhir kata menjadi oy.*/
	s = DiftongNormalize(s, `ai$`, "ay")
	s = DiftongNormalize(s, `au$`, "aw")
	s = DiftongNormalize(s, `oi$`, "oy")

	/* 4. Normalkan semivokal: ubah konsonan-y menjadi konsonan-i,
	iy menjadi i dan uw menjadi u.*/
	s = SemiVocalNormalize(s)
	s = strings.ReplaceAll(s, "iy", "i")
	s = strings.ReplaceAll(s, "uw", "u")

	/* 5. Normalkan konsonan yang berbunyi nyaris sama:
	   ubah kh dan q menjadi k, sy menjadi s, v menjadi f, z menjadi j,
	   d menjadi t, b menjadi p (mungkin masih ada yang kurang atau salah).*/
	s = strings.ReplaceAll(s, "kh", "k")
	s = strings.ReplaceAll(s, "q", "k")
	s = strings.ReplaceAll(s, "sy", "s")
	s = strings.ReplaceAll(s, "v", "f")
	s = strings.ReplaceAll(s, "z", "j")
	s = strings.ReplaceAll(s, "d", "t")
	s = strings.ReplaceAll(s, "b", "p")

	// 6. Normalkan ‘x’: ubah x menjadi ks
	s = strings.ReplaceAll(s, "x", "ks")

	/* 7. Ubah konsonan compound yang tersisa menjadi satu karakter:
	   ng menjadi d dan ny menjadi b.*/
	s = strings.ReplaceAll(s, "ng", "d")
	s = strings.ReplaceAll(s, "ny", "b")

	// 8. Normalkan h diam: ubah konsonan-h-vokal menjadi konsonan-vokal saja.
	s = DeadHRemoval(s)

	// 9. Hapus semua huruf vokal.
	s = DiftongNormalize(s, `[aiueo]+`, "")
	return s
}

func RemoveConsecutiveConsonants(s string) string {
	re := regexp.MustCompile(`c{2,}|b{2,}|d{2,}|g{2,}|f{2,}|h{2,}|k{2,}|j{2,}|m{2,}|l{2,}|n{2,}|q{2,}|p{2,}|s{2,}|r{2,}|t{2,}|w{2,}|v{2,}|y{2,}|x{2,}|z{2,}`)
	same := re.FindAllString(s, -1)
	for _, val := range same {

		consonant := string([]rune(val)[:1])
		s = strings.ReplaceAll(s, val, consonant)
	}
	return s
}

func DiftongNormalize(s, regex, new string) string {
	re := regexp.MustCompile(regex)
	same := re.FindAllString(s, -1)
	for _, val := range same {
		s = strings.ReplaceAll(s, val, new)
	}
	return s
}

func SemiVocalNormalize(s string) string {

	re := regexp.MustCompile(`cy{1,}|by{1,}|dy{1,}|gy{1,}|fy{1,}|hy{1,}|ky{1,}|jy{1,}|my{1,}|ly{1,}|ny{1,}|qy{1,}|py{1,}|sy{1,}|ry{1,}|ty{1,}|wy{1,}|vy{1,}|yy{1,}|xy{1,}|zy{1,}`)
	same := re.FindAllString(s, -1)
	for _, val := range same {
		consonant := string([]rune(val)[:1])
		//log.Print(consonant)
		s = strings.ReplaceAll(s, val, consonant+"i")
	}
	return s
}

func DeadHRemoval(s string) string {
	for i := 0; i < len(s); i++ {
		if i+3 < len(s) {
			if syllable := string([]rune(s)[i : i+3]); strings.ContainsAny(syllable, "h") && len(syllable) == 3 {
				hIndex := strings.IndexAny(syllable, "h")
				if hIndex-1 == -1 || hIndex+1 == 3 {
				} else {
					consonant := string([]rune(syllable)[hIndex-1])
					vocal := string([]rune(syllable)[hIndex+1])
					if strings.ContainsAny(consonant, "cbdgfhkjmlnqpsrtwvyxz") && strings.ContainsAny(vocal, "aiueo") {
						new := strings.ReplaceAll(syllable, "h", "")
						s = strings.ReplaceAll(s, syllable, new)
					}
				}
			}
		}
	}
	return s
}
