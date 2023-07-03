// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package render

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"
	"github.com/googlecodelabs/tools/claat/nodes"
)

// MD renders nodes as markdown for the target env.
func MD(ctx Context, nodes ...nodes.Node) (string, error) {
	var buf bytes.Buffer
	if err := WriteMD(&buf, ctx.Env, ctx.Format, nodes...); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// WriteMD does the same as MD but outputs rendered markup to w.
func WriteMD(w io.Writer, env string, fmt string, nodes ...nodes.Node) error {
	mw := mdWriter{w: w, env: env, format: fmt, Prefix: []byte("")}
	return mw.write(nodes...)
}

type mdWriter struct {
	w                  io.Writer // output writer
	env                string    // target environment
	format             string    // target template
	err                error     // error during any writeXxx methods
	lineStart          bool
	isWritingTableCell bool   // used to override lineStart for correct cell formatting
	isWritingList      bool   // used for override newblock when needed
	Prefix             []byte // prefix for e.g. blockquote content
}

func (mw *mdWriter) writeBytes(b []byte) {
	if mw.err != nil {
		return
	}
	if mw.lineStart {
		_, mw.err = mw.w.Write(mw.Prefix)
	}
	mw.lineStart = len(b) > 0 && b[len(b)-1] == '\n'
	_, mw.err = mw.w.Write(b)
}

func (mw *mdWriter) writeString(s string) {
	mw.writeBytes([]byte(s))
}

func (mw *mdWriter) writeEscape(s string) {
	s = html.EscapeString(s)
	mw.writeString(ReplaceDoubleCurlyBracketsWithEntity(s))
}

func (mw *mdWriter) space() {
	if !mw.lineStart {
		mw.writeString(" ")
	}
}

func (mw *mdWriter) newBlock() {
	if !mw.lineStart {
		mw.writeString("\n")
		mw.writeString("\n")
	}

  // Todo: Add line breaks in Text block 
  if mw.text != nil {
	   mw.writeString("\n")
//    n := mw.TextNode
//    tr := strings.TrimLeft(n.Value, "\t\n\r\f\v")
//
//    if strings.Contains(tr, "</ql-multiple-choice-probe>"){
//	    mw.writeString("\n")
//    } else {
//	    mw.writeString("\n\n")
//    }
  }
}

func (mw *mdWriter) matchEnv(v []string) bool {
	if len(v) == 0 || mw.env == "" {
		return true
	}
	i := sort.SearchStrings(v, mw.env)
	return i < len(v) && v[i] == mw.env
}

func (mw *mdWriter) write(nodesToWrite ...nodes.Node) error {
	for _, n := range nodesToWrite {
		if !mw.matchEnv(n.Env()) {
			continue
		}
		switch n := n.(type) {
		case *nodes.TextNode:
			mw.text(n)
		case *nodes.ImageNode:
			mw.image(n)
		case *nodes.URLNode:
			mw.url(n)
		case *nodes.ButtonNode:
			mw.write(n.Content.Nodes...)
		case *nodes.CodeNode:
			mw.code(n)
		case *nodes.ListNode:
			mw.list(n)
		case *nodes.ImportNode:
			if len(n.Content.Nodes) == 0 {
				break
			}
			mw.write(n.Content.Nodes...)
      // IMPORTS are handled as TEXT
	    mw.writeString("\n")
		case *nodes.ItemsListNode:
			mw.itemsList(n)
		case *nodes.GridNode:
			mw.table(n)
		case *nodes.InfoboxNode:
			mw.infobox(n)
		case *nodes.SurveyNode:
			mw.survey(n)
		case *nodes.HeaderNode:
			mw.header(n)
		case *nodes.YouTubeNode:
			mw.youtube(n)
		}
		if mw.err != nil {
			return mw.err
		}
	}
	return nil
}

func (mw *mdWriter) text(n *nodes.TextNode) {
	// tr := strings.TrimLeft(n.Value, " \t\n\r\f\v")
	tr := strings.TrimLeft(n.Value, " \t\r\f\v")
	left := n.Value[0:(len(n.Value) - len(tr))]
	// t := strings.TrimRight(tr, " \t\n\r\f\v")
	t := strings.TrimRight(tr, " \t\r\f\v")
	right := tr[len(t):len(tr)]

  if strings.Contains(t, "\n\n"){
	  mw.writeString("\n\n")
  }

  // TODO: Add line break before IMPORT
  if strings.Contains(t, "[["){
//	  mw.writeString("\n\n")
//	  mw.writeString("IMPORT statement")
	  mw.writeString("\n\n")
  }

  // TODO: Ensure break added for watermark
  var isEndWatermark bool

  // TODO: Ensure break added for watermark
  if strings.Contains(t, "Last Updated"){
	  mw.writeString("\n\n")
    isEndWatermark = true;
  }

  if strings.Contains(t, "Last Tested"){
	  mw.writeString("\n\n")
    isEndWatermark = true;
  }

   // TODO: Add a line break for paragraph
   if strings.Contains(t, "\n"){
 	  mw.writeString("\n\n")
   }

  // TODO: Automate Date update for String t


	mw.writeString(left)

  // TODO: Replace with HTML Code
	if n.Bold {
		//mw.writeString("**")
		mw.writeString("<strong>")
	}
  
	if n.Italic {
		mw.writeString("*")
	}

	if n.Code {
		mw.writeString("`")
	}

  // TODO: Enable ql-* custom block to be used in GDoc
//	t = strings.Replace(t, "<", "&lt;", -1)
//	t = strings.Replace(t, ">", "&gt;", -1)

  mw.writeString(t)

  if strings.Contains(t, "</ql-hint>"){
	  mw.writeString("\n\n")
  }

  if strings.Contains(t, "</ql-multiple-choice-probe>"){
	  // mw.writeString("\n\n")
	  mw.writeString("\n")
  }

  // TODO: Add line break after IMPORT
  if strings.Contains(t, "]]"){
//	  mw.writeString("\n\n")
//	  mw.writeString("IMPORT statement")
	  mw.writeString("\n\n")
  }

  // TODO: Ensure break added for watermark
  // Ensure line break added for watermark
  if isEndWatermark {
	  mw.writeString("\n\n")
  }
  

  // TODO: Replace HTML formatting with HTML Code
	if n.Code {
		mw.writeString("`")
	}

	if n.Italic {
		mw.writeString("*")
	}

	if n.Bold {
    // TODO: Amend so strong is rendered within the document
    // Use strong as this can be rendered inside an Info/Warning box
		//mw.writeString("**")
		mw.writeString("</strong>")
	}

	mw.writeString(right)
  // TODO: Not Codelab compatible
//  if n.Bold {
//	  mw.writeString("\n\n")
//  }
}

func (mw *mdWriter) image(n *nodes.ImageNode) {
	mw.space()
	mw.writeString("<img ")
	mw.writeString(fmt.Sprintf("src=%q ", n.Src))

	if n.Alt != "" {
		mw.writeString(fmt.Sprintf("alt=%q ", n.Alt))
	} else {
		mw.writeString(fmt.Sprintf("alt=%q ", path.Base(n.Src)))
	}

	if n.Title != "" {
		mw.writeString(fmt.Sprintf("title=%q ", n.Title))
	}

	// If available append width to the src string of the image.
	if n.Width > 0 {
		mw.writeString(fmt.Sprintf(" width=\"%.2f\" ", n.Width))
	}

	mw.writeString("/>")
	mw.writeString("\n")
	mw.writeString("\n")
}

func (mw *mdWriter) url(n *nodes.URLNode) {
//	mw.space()
	if n.URL != "" {
		// Look-ahead for button syntax.
		if _, ok := n.Content.Nodes[0].(*nodes.ButtonNode); ok {
			mw.writeString("<button>")
		}
		mw.writeString("[")
	}
	mw.write(n.Content.Nodes...)
	if n.URL != "" {
		// escape parentheses
		strings.Replace(n.URL, "(", "%28", -1)
		strings.Replace(n.URL, ")", "%29", -1)
		mw.writeString("](")
		mw.writeString(n.URL)
		mw.writeString(")")
		if _, ok := n.Content.Nodes[0].(*nodes.ButtonNode); ok {
			// Look-ahead for button syntax.
			mw.writeString("</button>")
		}
	}
}

func (mw *mdWriter) code(n *nodes.CodeNode) {
  // Allow user defined code block
  writeCodeBlock := true 

	if n.Empty() {
		return
	}
	mw.newBlock()
	defer mw.writeString("\n")

	if n.Term {
    // User defined: Handle code ticks 
    if strings.Contains(n.Value, "```"){
      // Code block defined
      writeCodeBlock = false 
    } 

    // User defined: Handle ql-code-block 
    if strings.Contains(n.Value, "ql-code-block") {
      // Code block defined
      writeCodeBlock = false 
    } 

    // Default: Handle code block 
    if  writeCodeBlock {
      // Code block default
      writeCodeBlock = true 

	    mw.writeString("\n")
	    mw.writeString("<ql-code-block bash templated noWrap>")
    }
	} else {
		mw.writeString(n.Lang)
	}

	mw.writeString("\n")
	mw.writeString(n.Value)
  
	if !mw.lineStart {
		mw.writeString("\n")
	}

  // TODO: Write the closing code block 
  if (writeCodeBlock) {
	  mw.writeString("</ql-code-block>")
  }
	mw.writeString("\n")
}

func (mw *mdWriter) list(n *nodes.ListNode) {
	if n.Block() == true {
		mw.newBlock()
	}
	mw.write(n.Nodes...)
	if !mw.lineStart && !mw.isWritingTableCell {
		mw.writeString("\n")
	}
}

func (mw *mdWriter) itemsList(n *nodes.ItemsListNode) {
	mw.isWritingList = true
	if n.Block() == true {
		mw.newBlock()
	}

  // TODO: Add line break before list
  mw.writeString("\n")

  // TODO: Replace with HTML Unordered List

	for i, item := range n.Items {
		s := "* "
		if n.Type() == nodes.NodeItemsList && n.Start > 0 {
			s = strconv.Itoa(i+n.Start) + ". "
		}
		mw.writeString(s)
		mw.write(item.Nodes...)
		if !mw.lineStart {
			mw.writeString("\n")
		}
	}
  // TODO: Add list space
  mw.writeString("\n")
	mw.isWritingList = false
}

func (mw *mdWriter) infobox(n *nodes.InfoboxNode) {
	// InfoBoxes are comprised of a ListNode with the contents of the InfoBox.
	// Writing the ListNode directly results in extra newlines in the md output
	// which breaks the formatting. So instead, write the ListNode's children
	// directly and don't write the ListNode itself.
	mw.newBlock()
  // TODO: Replace aside with infobox/warningbox
	// k := "aside positive"
	k := "<ql-infobox>"
	if n.Kind == nodes.InfoboxNegative {
		// k = "aside negative"
		k = "<ql-warningbox>"
	}
	mw.Prefix = []byte("")
	mw.writeString(k)
	mw.writeString("\n")

//	t = strings.Replace(t, "<", "&lt;", -1)
	for _, cn := range n.Content.Nodes {
		mw.write(cn)
	}

  // TODO: Close 
	mw.Prefix = []byte("")

  // TODO: Cloud the info/warningbox
	// k := "aside positive"
	k = "</ql-infobox>"
	if n.Kind == nodes.InfoboxNegative {
		// k = "aside negative"
		k = "</ql-warningbox>"
	}
	mw.Prefix = []byte("")
	mw.writeString(k)
	mw.writeString("\n")
}

func (mw *mdWriter) survey(n *nodes.SurveyNode) {
	mw.newBlock()
	mw.writeString("<form>")
	mw.writeString("\n")
	for _, g := range n.Groups {
		mw.writeString("<name>")
		mw.writeEscape(g.Name)
		mw.writeString("</name>")
		mw.writeString("\n")
		for _, o := range g.Options {
			mw.writeString("<input value=\"")
			mw.writeEscape(o)
			mw.writeString("\">")
			mw.writeString("\n")
		}
	}
	mw.writeString("</form>")
}

func (mw *mdWriter) header(n *nodes.HeaderNode) {
	mw.newBlock()
	mw.writeString(strings.Repeat("#", n.Level+1))
	mw.writeString(" ")
	mw.write(n.Content.Nodes...)
	if !mw.lineStart {
		mw.writeString("\n")
	}
}

func (mw *mdWriter) youtube(n *nodes.YouTubeNode) {
	if !mw.isWritingList {
		mw.newBlock()
	}

  // TODO: Video should be on a new Block
	mw.newBlock()	

	mw.writeString("\n")
  // TODO: Replace video control with ql-video element
	// mw.writeString(fmt.Sprintf(`<video id="%s"></video>`, n.VideoID))
	mw.writeString(fmt.Sprintf(`<ql-video youtubeId="%s"></ql-video>`, n.VideoID))
}

func (mw *mdWriter) table(n *nodes.GridNode) {
	// If table content is empty, don't output the table.
	if n.Empty() {
		return
	}

	mw.writeString("\n")
	maxcols := maxColsInTable(n)
	for rowIndex, row := range n.Rows {
		mw.writeString("|")
		for _, cell := range row {
			mw.isWritingTableCell = true
			mw.writeString(" ")

			// Check cell content for newlines and replace with inline HTML if newlines are present.
			var nw bytes.Buffer
			WriteMD(&nw, mw.env, mw.format, cell.Content.Nodes...)
			if bytes.ContainsRune(nw.Bytes(), '\n') {
				for _, cn := range cell.Content.Nodes {
					cn.MutateBlock(false) // don't treat content as a new block
					var nw2 bytes.Buffer
					WriteHTML(&nw2, mw.env, mw.format, cn)
					mw.writeBytes(bytes.Replace(nw2.Bytes(), []byte("\n"), []byte(""), -1))
				}
			} else {
				mw.writeBytes(nw.Bytes())
			}

			mw.writeString(" |")
		}
		if rowIndex == 0 && len(row) < maxcols {
			for i := 0; i < maxcols-len(row); i++ {
				mw.writeString(" |")
			}
		}
		mw.writeString("\n")

		// Write header bottom border
		if rowIndex == 0 {
			mw.writeString("|")
			for i := 0; i < maxcols; i++ {
				mw.writeString(" --- |")
			}
			mw.writeString("\n")
		}

		mw.isWritingTableCell = false
	}
}

func maxColsInTable(n *nodes.GridNode) int {
	m := 0
	for _, row := range n.Rows {
		if len(row) > m {
			m = len(row)
		}
	}
	return m
}
