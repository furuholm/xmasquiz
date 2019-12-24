package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell"

	"github.com/rivo/tview"
)

type question struct {
	text    string
	answers []string
	hints   map[string]string
}

var questions = []question{
	{
		// text: `God Jul!`,
		text: `God Jul!

Välkommen till årets skattjakt.

Som ni känner till så är tomten utbränd. Därför har jag, Siri, fått uppdraget att hålla i skattjakten.

Er uppgift är att leta reda på frågor som finns gömda i huset. Om ni svarar rätt får ni gå vidare.

Om ni inte kommer på svaret kan ni få en ledtråd av någon vuxen.

Ska vi börja?`,
		answers: []string{"ja"},
		hints: map[string]string{
			"nej": "vadå, vill ni inte ha nån jul i år eller? Försök igen.",
		},
	},
	{
		text: `Bra! Då kör vi!

Fråga 1 hittar ni i en apparat som man använder för att grilla smörgåsar.

Leta reda på frågan. Skriv sedan svaret på frågan nedanför.`,
		answers: []string{"pennywise"},
	},
	{
		text: `Mycket riktigt! Och kom ihåg barn att man aldrig ska bada med sin brödrost!

Fråga 2 hittar ni på ett ställe som inte är utomhus men ändå är kallt och som kan värmas med eld.

Leta reda på frågan. Skriv sedan svaret på frågan nedanför.`,
		answers: []string{"italien"},
	},
	{
		text: `Utmärkt! Om ni vill veta mer om Italien så vet ni nog vem ni ska fråga.

Hursomhelst, fråga 3 hittar ni där bönor blir till pulver.

Leta reda på frågan. Skriv sedan svaret på frågan nedanför.`,
		answers: []string{"grimma"},
	},
	{
		text: `Helt rätt! Det här går ju fint! Hoppas ni gillar mina kluriga frågor!

Fråga 4 hittar ni vid en apparat som kan upptäcka rök.`,
		answers: []string{"robin hood"},
	},
	{
		text: `Korrekt! Nu börjar det närma sig...
		
Den sista frågan hittar ni om ni löser den här gåtan: Bosse vill inte sova på ett torg, han trivs bättre i sin...

Ja, var trivs han?

Leta reda på frågan. Skriv sedan svaret på frågan nedanför.`,
		answers: []string{"4"},
	},
	{
		text: `Alla frågor avklarade!!!
		
Det var ett nöje att ha skattjakt med er! Men ni måste lösa en sista gåta för att hitta ert mål:

Här kommer den: skatten hittar ni under vägen upp...`,
	},
}

func main() {
	app := tview.NewApplication()
	index := 0

	textView := tview.NewTextView().
		SetRegions(true).
		SetWordWrap(true)
	textGrid := tview.NewGrid().
		SetRows(15, 10).
		AddItem(textView, 0, 0, 1, 1, 10, 0, false)
	inputField := tview.NewInputField().
		SetLabel("Svar:")
	grid := tview.NewGrid().
		SetRows(15, 10).
		AddItem(textView, 0, 0, 1, 1, 10, 0, false).
		AddItem(inputField, 1, 0, 1, 1, 0, 0, true)

	textFrame := tview.NewFrame(textView).
		SetBorders(10, 2, 2, 2, 4, 4)

	setQuestion := func(q question) {
		textView.Clear()
		text := q.text
		fmt.Fprintf(textView, "%s ", text)
		if len(q.answers) == 0 {
			app.SetRoot(textGrid, true)
			app.ForceDraw()
			say(q.text)
		} else {
			app.SetRoot(textGrid, true)
			app.ForceDraw()
			say(q.text)
			app.SetRoot(grid, true)
		}
	}

	next := func(key tcell.Key) {
		answer := inputField.GetText()
		if answer == "" {
			return
		}
		if index+1 < len(questions) {
			q := questions[index]
			normalizedAnswer := strings.ToLower(strings.Trim(answer, " !."))
			if exists(normalizedAnswer, q.answers) {
				// correct answer
				index += 1
				setQuestion(questions[index])
			} else if hint, ok := q.hints[normalizedAnswer]; ok {
				say(hint)
			} else {
				say("Fel svar, försök igen...")
			}
			inputField.SetText("")
		}
	}

	textView.
		SetTextAlign(tview.AlignCenter).
		SetDoneFunc(func(key tcell.Key) {
			q := questions[0]
			setQuestion(q)
		}).
		SetText("Tryck enter för att börja!")

	inputField.SetDoneFunc(next)

	if err := app.SetRoot(textFrame, true).Run(); err != nil {
		panic(err)
	}
}

func say(text string) {
	cmd := exec.Command("say", "-v", "alva", text)
	_ = cmd.Run()
}

func exists(needle string, haystack []string) bool {
	for _, s := range haystack {
		if needle == s {
			return true
		}
	}
	return false
}
