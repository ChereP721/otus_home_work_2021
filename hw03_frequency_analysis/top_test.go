package hw03frequencyanalysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("small count words", func(t *testing.T) {
		text := `
1
2 2
3 3 3
4 4 4 4
		`
		expected := []string{
			"4",
			"3",
			"2",
			"1",
		}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("int overage", func(t *testing.T) {
		text := strings.Repeat(" int overage", 128) + " a one word"
		expected := []string{
			"a",
			"one",
			"word",
			"int",
			"overage",
		}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("unicode", func(t *testing.T) {
		text := `
\u0011
\u0012 \u0012
\u0013 \u0013
\u0014 \u0014 \u0014
		`
		expected := []string{
			`\u0014`,
			`\u0012`,
			`\u0013`,
			`\u0011`,
		}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("hard language and emoji", func(t *testing.T) {
		text := `
你
好 好
第 第 第 第
三 三 三 三
🤯 🤯 🤯 🤯
😀 😀 😀 😀
🤗 🤗 🤗 🤗 🤗
☠ ☠ ☠ ☠ ☠
		`
		expected := []string{
			"☠",
			"🤗",
			"三",
			"第",
			"😀",
			"🤯",
			"好",
			"你",
		}
		require.Equal(t, expected, Top10(text))
	})

	taskWithAsteriskIsCompleted = false
	t.Run("positive test", func(t *testing.T) {
		expected := []string{
			"он",        // 8
			"а",         // 6
			"и",         // 6
			"ты",        // 5
			"что",       // 5
			"-",         // 4
			"Кристофер", // 4
			"если",      // 4
			"не",        // 4
			"то",        // 4
		}
		require.Equal(t, expected, Top10(text))
	},
	)

	taskWithAsteriskIsCompleted = true
	t.Run("positive test with asterisk", func(t *testing.T) {
		expected := []string{
			"а",         // 8
			"он",        // 8
			"и",         // 6
			"ты",        // 5
			"что",       // 5
			"в",         // 4
			"его",       // 4
			"если",      // 4
			"кристофер", // 4
			"не",        // 4
		}
		require.Equal(t, expected, Top10(text))
	},
	)
}

func TestNormalizeWord(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   string
		expectedWA string
	}{
		{name: "simple word", input: "Как", expected: "Как", expectedWA: "как"},
		{name: "dash word", input: " - ", expected: "-", expectedWA: ""},
		{name: "punctuation marks", input: "видите,", expected: "видите,", expectedWA: "видите"},
		{name: "word with hyphen", input: "бум-бум-бум", expected: "бум-бум-бум", expectedWA: "бум-бум-бум"},
		{name: "word with hyphen at end", input: "бум-бум-бум-", expected: "бум-бум-бум-", expectedWA: "бум-бум-бум"},
	}
	for _, tc := range tests {
		tc := tc
		taskWithAsteriskIsCompleted = false
		t.Run(tc.name+" ("+tc.input+")", func(t *testing.T) {
			result := normalizeWord(tc.input)
			require.Equal(t, tc.expected, result)
		})
		taskWithAsteriskIsCompleted = true
		t.Run(tc.name+" with asterisk ("+tc.input+")", func(t *testing.T) {
			result := normalizeWord(tc.input)
			require.Equal(t, tc.expectedWA, result)
		})
	}
}
