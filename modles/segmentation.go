package modles

import (
	"apiTools/libs/config"
	"github.com/pkg/errors"
)

// 分词Api

const (
	segTextToLong = "segmentation text is to long"
)

type SegForm struct {
	Text string `form:"text" json:"text" xml:"text" binding:"required"`
}

type SegResponse struct {
	Result []string `json:"result"`
}

func SegTextCut(form *SegForm) (resp *SegResponse, msg string, err error) {
	if len(form.Text) > 100 {
		err = errors.New(segTextToLong)
		msg = segTextToLong
	}
	cutSlice := config.Seg.Cut(form.Text, true)
	resp = &SegResponse{
		Result: cutSlice,
	}
	return
}
