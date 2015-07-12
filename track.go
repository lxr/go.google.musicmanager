package musicmanager

import (
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/dhowden/tag"

	"github.com/golang/protobuf/proto"
	pb "google-musicmanager-go/proto"
)

// readMetadataErr is logged if ReadMetadataFromMP3 does not know how
// to convert a tag's contents into a string; this is not a fatal error,
// and causes a tag to be ignored instead.
const readMetadataErr = "readMetadataFromMP3 does not know how to " +
	"convert tag %s of type %T into string, ignoring " +
	"(github.com/dhowden/tag has likely been updated)"

// ReadMetadataFromMP3 extracts metadata from the given MP3 file.
// Its behavior is undefined if the argument is not an MP3 file.
// It is called automatically by TracksInsertCall.Do if no custom
// metadata was provided by the caller.
func ReadMetadataFromMP3(f io.ReadSeeker) (*pb.Track, error) {
	info, err := mp3Fexamine(f, false)
	if err = rewind(f, err); err != nil {
		return nil, err
	}
	sum, err := tag.Sum(f)
	if err = rewind(f, err); err != nil {
		return nil, err
	}
	metadata, err := tag.ReadFrom(f)
	if err == tag.ErrNoTagsFound {
		err = nil
	}
	if err = rewind(f, err); err != nil {
		return nil, err
	}
	var tags map[string]interface{}
	if metadata != nil {
		tags = metadata.Raw()
	}
	track := &pb.Track{
		Title:           proto.String("Untitled Track"),
		ClientId:        proto.String(sum),
		DurationMillis:  proto.Int64(int64(info.Length * 1000)),
		OriginalBitRate: proto.Int32(int32(info.Bitrate)),
		ContentType:     pb.Track_MP3.Enum(),
		EstimatedSize:   proto.Int64(info.Size),
	}
	ptrack := reflect.ValueOf(track).Elem()
	for field, record := range fieldfromframe {
		if value, ok := tags[record.Frame]; ok {
			var str string
			switch value := value.(type) {
			case string:
				str = value
			case *tag.Comm:
				str = value.Text
			default:
				log.Printf(readMetadataErr, record.Frame, value)
				continue
			}
			res := record.ConvertFunc(str)
			if res.IsValid() {
				ptrack.FieldByName(field).Set(res)
			}
		}
	}
	return track, nil
}

// rewind resets a Seeker to its beginning, unless given a non-nil
// error to propagate.
func rewind(f io.Seeker, err error) error {
	if err == nil {
		_, err = f.Seek(0, 0)
	}
	return err
}

// fieldfromframe describes how to map id3v2 tags to pb.Track fields.
var fieldfromframe = map[string]struct {
	Frame       string
	ConvertFunc func(s string) reflect.Value
}{
	"Title":           {"TIT2", id},
	"Artist":          {"TPE1", id},
	"Composer":        {"TCOM", id},
	"Album":           {"TALB", id},
	"AlbumArtist":     {"TPE2", id},
	"Year":            {"TYER", atoi},
	"Comment":         {"COMM", id},
	"TrackNumber":     {"TRCK", splitter(0)},
	"Genre":           {"TCON", id},
	"BeatsPerMinute":  {"TBPM", atoi},
	"PlayCount":       {"PCNT", atoi},
	"TotalTrackCount": {"TRCK", splitter(1)},
	"DiscNumber":      {"TPOS", splitter(0)},
	"TotalDiscCount":  {"TPOS", splitter(1)},
	"Compilation":     {"TCMP", flag},
	// I don't know how dhowden/tag reports POPM
	//"Rating":        {"POPM", id},
}

// id returns the value of the tag as a string.
func id(s string) reflect.Value {
	return reflect.ValueOf(proto.String(s))
}

// atoi attempts to interpret the value of the tag as an interger,
// returning the zero Value on failure.
func atoi(s string) reflect.Value {
	i, err := strconv.Atoi(s)
	if err != nil {
		return reflect.Value{}
	}
	return reflect.ValueOf(proto.Int32(int32(i)))
}

// splitter returns a function that attempts to interpret either the
// first or second half of the tag (arguments 0 and 1 respectively,
// where the halves are separated by the character '/') as an integer.
func splitter(i int) func(s string) reflect.Value {
	return func(s string) reflect.Value {
		a := strings.SplitN(s, "/", 2)
		if len(a) < 2 {
			a = append(a, "")
		}
		return atoi(a[i])
	}
}

// flag returns true if the tag exists.
func flag(s string) reflect.Value {
	return reflect.ValueOf(proto.Bool(true))
}
