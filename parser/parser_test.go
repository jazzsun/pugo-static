package parser_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-xiaohei/pugo-static/model"
	"github.com/go-xiaohei/pugo-static/parser"
	. "github.com/smartystreets/goconvey/convey"
)

type demoData struct {
	Name     string
	Age      int
	IsFemale bool
}

var (
	p  = parser.NewCommonParser()
	p2 = parser.NewMdParser()
)

func TestIniBlock(t *testing.T) {
	Convey("test ini block", t, func() {
		Convey("test ini empty data", func() {
			ib := new(parser.IniBlock)
			ib.Write(nil)
			demo := new(demoData)
			err := ib.MapTo("", demo)
			So(err, ShouldBeNil)
			So(demo.Name, ShouldBeEmpty)

			ib.Write(nil)
			data := ib.MapHash("")
			So(data, ShouldBeEmpty)

			ib.Write(nil)
			data2 := ib.Keys("")
			So(data2, ShouldBeEmpty)
		})

		Convey("test wrong ini data", func() {
			ib := new(parser.IniBlock)
			ib.Write([]byte("aaaaaaaa"))
			demo := new(demoData)
			err := ib.MapTo("", demo)
			So(err, ShouldNotBeNil)
			So(demo.Name, ShouldBeEmpty)

			ib.Write([]byte("bbbbbb"))
			data := ib.MapHash("")
			So(data, ShouldBeEmpty)

			ib.Write([]byte("ccccc"))
			data2 := ib.Keys("")
			So(data2, ShouldBeEmpty)

			ib.Write([]byte("dddd"))
			data3 := ib.Item("abc")
			So(data3, ShouldBeEmpty)
		})
	})
}

func TestParser(t *testing.T) {
	Convey("test parser", t, func() {
		Convey("empty content", func() {
			blocks, err := p.Parse(nil)
			So(blocks, ShouldBeNil)
			So(err, ShouldBeNil)

			blocks, err = p.Parse([]byte(""))
			So(blocks, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("is common parser", func() {
			flag := p.Is([]byte(parser.MD_PARSER_PREFIX))
			So(flag, ShouldBeFalse)

			flag = p.Is([]byte(parser.COMMON_PARSER_PREFIX))
			So(flag, ShouldBeTrue)

			flag = p2.Is([]byte(parser.MD_PARSER_PREFIX))
			So(flag, ShouldBeTrue)
		})

		Convey("detect block", func() {
			block := p.Detect([]byte("xxx"))
			So(block, ShouldBeNil)

			block = p2.Detect([]byte("xxx"))
			So(block, ShouldBeNil)

			block = p.Detect([]byte("ini"))
			So(block, ShouldNotBeNil)
			So(block.Type(), ShouldEqual, parser.BLOCK_INI)
		})

		Convey("first block error", func() {
			blocks, err := p.Parse([]byte("-----xxx\ncontent"))
			So(blocks, ShouldBeNil)
			So(err, ShouldNotBeNil)

			blocks, err = p2.Parse([]byte("```xxx\ncontent```"))
			So(blocks, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		blocks, err := p.Parse([]byte("\n\n-----ini\ncontent"))
		So(err, ShouldBeNil)
		So(blocks, ShouldNotBeNil)
	})
}

func TestParseMeta(t *testing.T) {
	Convey("parse meta", t, func() {
		bytes, err := ioutil.ReadFile("../source/meta.md")
		So(err, ShouldBeNil)
		blocks, err := p2.Parse(bytes)
		So(err, ShouldBeNil)

		Convey("check meta block", func() {
			So(blocks, ShouldHaveLength, 1)
			So(blocks[0].Type(), ShouldEqual, parser.BLOCK_INI)

			Convey("use meta block", func() {
				b, ok := blocks[0].(parser.MetaBlock)
				So(ok, ShouldBeTrue)
				So(b.Item("meta", "title"), ShouldEqual, "Pugo.Static")

				meta, _, _, _, err := model.NewAllMeta(blocks)
				So(err, ShouldBeNil)
				So(meta.Title, ShouldEqual, b.Item("meta", "title"))
			})
		})
	})
}

func TestPostMeta(t *testing.T) {
	Convey("parse post", t, func() {
		bytes, err := ioutil.ReadFile("../source/post/welcome.md")
		So(err, ShouldBeNil)
		blocks, err := p2.Parse(bytes)
		So(err, ShouldBeNil)

		Convey("check post blocks", func() {
			So(blocks, ShouldHaveLength, 2)
			So(blocks[0].Type(), ShouldEqual, parser.BLOCK_INI)
			So(blocks[1].Type(), ShouldEqual, parser.BLOCK_MARKDOWN)

			Convey("use post blocks", func() {
				b, ok := blocks[0].(parser.MetaBlock)
				So(ok, ShouldBeTrue)
				So(b.Item("title"), ShouldEqual, "Welcome")

				fi, _ := os.Stat("../source/post/welcome.md")
				post, err := model.NewPost(blocks, fi)
				So(err, ShouldBeNil)
				So(post.Title, ShouldEqual, b.Item("title"))
			})
		})
	})
}

func TestPageMeta(t *testing.T) {
	Convey("parse page with MdParser", t, func() {
		bytes, err := ioutil.ReadFile("../source/page/about.md")
		So(err, ShouldBeNil)
		blocks, err := p2.Parse(bytes)
		So(err, ShouldBeNil)
		So(blocks, ShouldHaveLength, 2)

		Convey("check page blocks", func() {
			So(blocks[0].Type(), ShouldEqual, parser.BLOCK_INI)
			So(blocks[1].Type(), ShouldEqual, parser.BLOCK_MARKDOWN)

			Convey("use page blocks", func() {
				b, ok := blocks[0].(parser.MetaBlock)
				So(ok, ShouldBeTrue)
				So(b.Item("title"), ShouldEqual, "About Pugo.Static")

				fi, _ := os.Stat("../source/page/about.md")
				page, err := model.NewPage(blocks, fi)
				So(err, ShouldBeNil)
				So(page.Title, ShouldEqual, b.Item("title"))
				So(page.Meta["metadata"], ShouldEqual, b.Item("meta", "metadata"))
			})
		})
	})
}
