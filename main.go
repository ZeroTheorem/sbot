package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ZeroTheorem/sbot/db"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	_ "modernc.org/sqlite"
)

type states struct {
	selectsYear  bool
	selectsMonth bool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	b, err := tele.NewBot(tele.Settings{
		Token:     os.Getenv("TOKEN"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
	})

	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open("sqlite", "file:mydb.db")

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	q := db.New(conn)

	var message *tele.Message

	// -- Section: create main menu
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnAdd := menu.Data("👌👈 +1", "add")
	btnMonthTotal := menu.Data("🔞💦 Итоги месяца! 💦🔞", "monthTotal")
	btnYearTotal := menu.Data("🔞💦 Итоги года! 💦🔞", "yearTotal")
	menu.Inline(
		menu.Row(btnAdd),
		menu.Row(btnMonthTotal),
		menu.Row(btnYearTotal),
	)
	// -- end section

	// -- Section: create inline keyboard
	yearSelector := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnCertainYearTotal := yearSelector.Data("Итоги 🔞 за другой год!", "certainYearTotal", "certainYearTotal")
	btnPrev := yearSelector.Data("Назад", "prev")
	yearSelector.Inline(
		yearSelector.Row(btnCertainYearTotal),
		yearSelector.Row(btnPrev),
	)
	monthSelector := &tele.ReplyMarkup{}
	btnCertainMonthTotal := yearSelector.Data("Итоги 🔞 за другой месяц!", "certainMonthTotal", "certainMonthTotal")
	monthSelector.Inline(
		yearSelector.Row(btnCertainMonthTotal),
		yearSelector.Row(btnPrev),
	)
	// -- end section

	// -- Section: initialize states
	st := states{}
	// -- end section

	// -- Section: initialize map month
	months := map[int64]string{
		1:  "Январь",
		2:  "Февраль",
		3:  "Март",
		4:  "Апрель",
		5:  "Май",
		6:  "Июнь",
		7:  "Июль",
		8:  "Август",
		9:  "Сентябрь",
		10: "Октябрь",
		11: "Ноябрь",
		12: "Декабрь",
	}
	// -- end section
	b.Handle("/menu", func(c tele.Context) error {
		m, err := b.Send(tele.ChatID(c.Chat().ID), "<b>Привет, я помогу вам с подсчетом ваших 'близостей'</b>😉", menu)
		if err != nil {
			return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
		}
		message = m
		return nil
	})
	b.Handle(&btnAdd, func(c tele.Context) error {
		now := time.Now()
		err := q.CreateRow(ctx, db.CreateRowParams{
			Month: int64(now.Month()),
			Year:  int64(now.Year()),
		})
		if err != nil {
			return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
		}
		count, err := q.GetAllByMonth(ctx, int64(now.Month()))
		if err != nil {
			return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
		}
		_, err = b.Edit(message, fmt.Sprintf("<b>В этом месяце, мои половые гиганты, вы уже перепехнулись 🔞</b>\n\n🔞💦👌👈:  <b>%v</b>", count), menu)
		return nil
	})
	b.Handle(&btnYearTotal, func(c tele.Context) error {
		yearTotal, err := q.GetAllByYear(ctx, int64(time.Now().Year()))
		if err != nil {
			return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
		}
		_, err = b.Edit(message, fmt.Sprintf("<b>За весь этот год джанджубанжу было неимоверное колличество раз!</b>\n\n🔞💦👌👈:  <b>%v</b>", yearTotal), yearSelector)
		return nil
	})

	b.Handle(&btnCertainYearTotal, func(c tele.Context) error {
		st.selectsYear = true
		return c.Send("<b>Введите год в формате: %YYYY (2025, 2026)</b>")
	})
	b.Handle(&btnMonthTotal, func(c tele.Context) error {
		monthTotal, err := q.GetAllByMonth(ctx, int64(time.Now().Month()))
		if err != nil {
			return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
		}
		_, err = b.Edit(message, fmt.Sprintf("<b>В этом месяце, кекса 🔞 у вас было больше чем у 99%% людей</b>\n\n🔞💦👌👈:  <b>%v</b>", monthTotal), monthSelector)
		return nil
	})
	b.Handle(&btnCertainMonthTotal, func(c tele.Context) error {
		st.selectsMonth = true
		return c.Send("<b>Введите номер месяца в формате: %M (1, 2, 3)</b>")
	})
	b.Handle("/delete", func(c tele.Context) error {
		err := q.DeleteLast(ctx)
		if err != nil {
			return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
		}
		return c.Send("<b>Записть успешно удалена</b>")
	})
	b.Handle(&btnPrev, func(c tele.Context) error {
		_, err = b.Edit(message, "<b>Привет, я помогу вам с подсчетом ваших 'близостей'</b>😉", menu)
		return nil
	})
	b.Handle(tele.OnText, func(c tele.Context) error {
		switch {
		case st.selectsYear:
			msg := c.Message().Text
			i, err := strconv.ParseInt(msg, 10, 64)
			if err != nil {
				return c.Send("<b>Не похоже что это чилсо</b>")
			}
			yearTotal, err := q.GetAllByYear(ctx, i)
			if err != nil {
				return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
			}
			st.selectsYear = false
			m, err := b.Send(tele.ChatID(c.Chat().ID), fmt.Sprintf("<b>За весь <i>%v</i> год у вас было столько 🔞, что остальным остается только завидывать</b>\n\n🔞💦👌👈:  <b>%v</b>", msg, yearTotal), menu)
			if err != nil {
				return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
			}
			message = m
			return nil
		case st.selectsMonth:
			msg := c.Message().Text
			i, err := strconv.ParseInt(msg, 10, 64)
			if err != nil {
				return c.Send("<b>Не похоже что это чилсо</b>")
			}
			if i < 1 || i > 12 {
				return c.Send("<b>Введите корректный месяц</b>")

			}
			monthTotal, err := q.GetAllByMonth(ctx, i)
			if err != nil {
				return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
			}
			st.selectsMonth = false
			m, err := b.Send(tele.ChatID(c.Chat().ID), fmt.Sprintf("<b>Всего за <i>%v</i> вы занялись '👌👈ЭТИМ👌👈' целых</b>\n\n🔞💦👌👈:  <b>%v</b>", months[i], monthTotal), menu)
			if err != nil {
				return c.Send(fmt.Sprintf("Oooops... something went wrong %v", err))
			}
			message = m
			return nil
		default:
			return nil
		}
	})
	b.Start()
}
