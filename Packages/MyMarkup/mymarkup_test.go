// mymarkup_test.go
// Tests for MyMarkup package
//
// 2025-07-05	PV 		First version

package MyMarkup

import (
	"testing"
)

func Test30(t *testing.T) {
	text := "Ceci est un texte à formater, avec un mot très long comme anticonstitutionnellement qui va être tronqué quand la largeur du rendu devient particulièrement petite."

    expected := `------------------------------
Ceci est un texte à formater, |
avec un mot très long comme   |
anticonstitutionnellement qui |
va être tronqué quand la      |
largeur du rendu devient      |
particulièrement petite.      |`

    s := BuildMarkupCore(text, true, 30)
    if s!=expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, s)
	}
}

func Test25(t *testing.T) {
    text := "Ceci est un texte à formater, avec un mot très long comme anticonstitutionnellement qui va être tronqué quand la largeur du rendu devient particulièrement petite."

    expected := `-------------------------
Ceci est un texte à      |
formater, avec un mot    |
très long comme          |
anticonstitutionnellement|
qui va être tronqué quand|
la largeur du rendu      |
devient particulièrement |
petite.                  |`

    s := BuildMarkupCore(text, true, 25)
    if s!=expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, s)
	}
}

func Test20(t *testing.T) {
    text := "Ceci est un texte à formater, avec un mot très long comme anticonstitutionnellement qui va être tronqué quand la largeur du rendu devient particulièrement petite."

    expected := `--------------------
Ceci est un texte à |
formater, avec un   |
mot très long comme |
anticonstitutionnell|
ement qui va être   |
tronqué quand la    |
largeur du rendu    |
devient             |
particulièrement    |
petite.             |`

    s := BuildMarkupCore(text, true, 20)
    if s!=expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, s)
	}
}

func Test15(t *testing.T) {
    text := "Ceci est un texte à formater, avec un mot très long comme anticonstitutionnellement qui va être tronqué quand la largeur du rendu devient particulièrement petite."

    expected := `---------------
Ceci est un    |
texte à        |
formater, avec |
un mot très    |
long comme     |
anticonstitutio|
nnellement qui |
va être tronqué|
quand la       |
largeur du     |
rendu devient  |
particulièremen|
t petite.      |`

    s := BuildMarkupCore(text, true, 15)
    if s!=expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, s)
	}
}

func Test10(t *testing.T) {
    text := "Ceci est un texte à formater, avec un mot très long comme anticonstitutionnellement qui va être tronqué quand la largeur du rendu devient particulièrement petite."

    expected := `----------
Ceci est  |
un texte à|
formater, |
avec un   |
mot très  |
long comme|
anticonsti|
tutionnell|
ement qui |
va être   |
tronqué   |
quand la  |
largeur du|
rendu     |
devient   |
particuliè|
rement    |
petite.   |`

    s := BuildMarkupCore(text, true, 10)
    if s!=expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, s)
	}
}