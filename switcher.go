package switcher

const (
	enLayout = "`1234567890-=\\qwertyuiop[]asdfghjkl;'zxcvbnm,./~!@#$%^&*()_+|QWERTYUIOP{}ASDFGHJKL:\"ZXCVBNM<>?"
	ruLayout = "ё1234567890-=\\йцукенгшщзхъфывапролджэячсмитьбю.Ё!\"№;%:?*()_+/ЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ,"
)

type Direction struct {
	Src *kLayout
	Tgt *kLayout
	Map map[rune]rune
}

type kLayout struct {
	Title  string
	Layout string
}

func newTranscodeDirection(a, b *kLayout) *Direction {
	dir := &Direction{
		Src: a,
		Tgt: b,
		Map: make(map[rune]rune),
	}
	rA, rB := []rune(a.Layout), []rune(b.Layout)
	for i, r := range rA {
		dir.Map[r] = rB[i]
	}
	return dir
}

type Transcoder struct {
	Directions []*Direction
}

func NewTranscoder() *Transcoder {
	en, ru := getLayouts()
	return &Transcoder{
		Directions: []*Direction{
			newTranscodeDirection(en, ru),
			newTranscodeDirection(ru, en),
		},
	}
}

func getLayouts() (A, B *kLayout) {
	return &kLayout{"EN", enLayout}, &kLayout{"RU", ruLayout}
}

func (d *Direction) CanTranscode(r rune) bool {
	_, ok := d.Map[r]
	return ok
}

func (d *Direction) Transcode(in []rune) (out []rune) {
	out = make([]rune, len(in))
	for i, r := range in {
		if r1, ok := d.Map[r]; ok {
			out[i] = r1
		} else {
			out[i] = r
		}
	}
	return
}

func (t *Transcoder) Transcode(str string) (res string) {
	if str == "" {
		return
	}

	var (
		word, out []rune
		dir       *Direction
	)

	for _, r := range str {
		runeDir, nDirs := (*Direction)(nil), 0
		for _, d := range t.Directions {
			if d.CanTranscode(r) {
				nDirs++
				runeDir = d
			}
		}
		if nDirs != 1 {
			runeDir = nil
		}

		if dir == nil || runeDir == nil || dir == runeDir {
			if dir == nil {
				dir = runeDir
			}
			word = append(word, r)
		} else {
			out = append(out, dir.Transcode(word)...)
			word = []rune{r}
			dir = runeDir
		}
	}

	if dir != nil {
		out = append(out, dir.Transcode(word)...)
	} else {
		out = append(out, word...)
	}
	res = string(out)
	return
}
