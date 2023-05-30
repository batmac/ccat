package mutators

import (
	"encoding/json"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/secretprovider"
)

const postURL = "https://translation.googleapis.com/language/translate/v2"

/*
Afrikaans 	af
Albanais 	sq
Amharique 	am
Arabe 	ar
Arménien 	hy
Assamais* 	as
Aymara* 	ay
Azéri 	az
Bambara* 	bm
Basque 	eu
Biélorusse 	be
Bengalî 	bn
Bhodjpouri* 	bho
Bosniaque 	bs
Bulgare 	bg
Catalan 	ca
Cebuano 	ceb
Chinois (simplifié) 	zh-CN ou zh (BCP-47)
Chinois (traditionnel) 	zh-TW (BCP-47)
Corse 	co
Croate 	hr
Tchèque 	cs
Danois 	da
Divéhi* 	dv
Dogri* 	doi
Néerlandais 	nl
Anglais 	en
Espéranto 	eo
Estonien 	et
Ewe* 	ee
Filipino (Tagalog) 	fil
Finnois 	fi
Français 	fr
Frison 	fy
Galicien 	gl
Géorgien 	ka
Allemand 	de
Grec 	el
Guarani* 	gn
Gujarâtî 	gu
Créole haïtien 	ht
Haoussa 	ha
Hawaïen 	haw
Hébreu 	he ou iw
Hindi 	hi
Hmong 	hmn
Hongrois 	hu
Islandais 	is
Igbo 	ig
Ilocano* 	ilo
Indonésien 	id
Irlandais 	ga
Italien 	it
Japonais 	ja
Javanais 	jv ou jw
Kannara 	kn
Kazakh 	kk
Khmer 	km
Kinyarwanda 	rw
Konkani* 	gom
Coréen 	ko
Krio* 	kri
Kurde 	ku
Kurde (Sorani)* 	ckb
Kirghyz 	ky
Laotien 	lo
Latin 	la
Letton 	lv
Lingala* 	ln
Lituanien 	lt
Luganda* 	lg
Luxembourgeois 	lb
Macédonien 	mk
Maïthili* 	mai
Malgache 	mg
Malais 	ms
Malayâlam 	ml
Maltais 	mt
Maori 	mi
Marathi 	mr
Meitei (Manipuri)* 	mni-Mtei
Mizo* 	lus
Mongol 	mn
Birman 	my
Népalais 	ne
Norvégien 	no
Nyanja (Chichewa) 	ny
Odia (Oriya) 	or
Oromo* 	om
Pachtô 	ps
Perse 	fa
Polonais 	pl
Portugais (Portugal, Brésil) 	pt
Panjabi 	pa
Quechua* 	qu
Roumain 	ro
Russe 	ru
Samoan 	sm
Sanskrit* 	sa
Gaélique (Écosse) 	gd
Sepedi* 	nso
Serbe 	sr
Sesotho 	st
Shona 	sn
Sindhî 	sd
Singhalais 	si
Slovaque 	sk
Slovène 	sl
Somali 	so
Spanish 	es
Soundanais 	su
Swahili 	sw
Suédois 	sv
Tagalog (philippin) 	tl
Tadjik 	tg
Tamoul 	ta
Tatar 	tt
Télougou 	te
Thaï 	th
Tigrinya* 	ti
Tsonga* 	ts
Turc 	tr
Turkmène 	tk
Twi (Akan)* 	ak
Ukrainien 	uk
Urdu 	ur
Ouïghour 	ug
Ouzbek 	uz
Vietnamien 	vi
Gallois 	cy
Xhosa 	xh
Yiddish 	yi
Yoruba 	yo
Zoulou 	zu
*/

func init() {
	defaultLanguage := "en"
	if t := os.Getenv("TARGET_LANGUAGE"); t != "" {
		defaultLanguage = t
	}
	singleRegister("translate", translate,
		withDescription("translate to X:en or $TARGET_LANGUAGE with google translate (needs a valid key in $GOOGLE_API_KEY)"),
		withConfigBuilder(stdConfigStringWithDefault(defaultLanguage)),
		withCategory("external APIs"),
	)
}

type Response struct {
	Data struct {
		Translations []struct {
			TranslatedText         string
			DetectedSourceLanguage string
		}
	}
}

func translate(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	targetLanguage := conf.(string)
	msg, err := io.ReadAll(r) // NOT streamable
	if err != nil {
		return 0, err
	}

	if len(msg) == 0 {
		return 0, nil
	}

	result := strings.Builder{}

	key, _ := secretprovider.GetSecret("translate", "GOOGLE_API_KEY")
	if key == "" {
		log.Fatal("no key found in $GOOGLE_API_KEY")
	}

	v := url.Values{}
	v.Set("key", key)
	v.Set("q", string(msg))
	v.Set("target", targetLanguage)

	res, err := http.PostForm(postURL, v)
	if err != nil {
		return 0, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	// fmt.Println(string(data))

	var d Response
	if err := json.Unmarshal(data, &d); err != nil {
		return 0, err
	}
	if len(d.Data.Translations) == 0 {
		log.Printf("didn't get any translation, got %s\n", string(data))
		return 0, nil
	}
	log.Debugf("Found a translation from language %s to %s\n", d.Data.Translations[0].DetectedSourceLanguage, targetLanguage)
	result.WriteString(html.UnescapeString(d.Data.Translations[0].TranslatedText))
	for _, txt := range d.Data.Translations[1:] {
		log.Debugf("Found an extra translation from language %s: %s\n", txt.DetectedSourceLanguage, txt.TranslatedText)
	}

	return io.Copy(w, strings.NewReader(result.String()))
}
