package templating_test

import (
	"github.com/brainboxweb/go-youtube-admin/templating"
	"strings"
	"testing"
)

func TestMusic(t *testing.T) {

	music := []string{"music1", "music2", "music3"}
	data := templating.YouTubeData{
		Music: music,
	}
	parsed, _ := templating.GetYouTubeBody(data, "youtube.txt")

	expected := "Music: music1 music2 music3"

	if !strings.Contains(parsed, expected) {
		t.Errorf("Expected music to be %s", music)
	}
}

//
//func TestTranscript(t *testing.T) {
//
//	id := "ididididid"
//	body := `the Body first line
//
//the body second line
//
//the body third line`
//	transcript := longTranscript
//	topResult := "http://number-one-on-google.com"
//	music := []string{"music1", "music2", "music3"}
//
//	data := templating.YouTubeData{
//		Id:         id,
//		Body:       body,
//		Transcript: transcript,
//		TopResult:  topResult,
//		Music:      music,
//	}
//
//	parsed := templating.GetYouTubeBody(data, "youtube.txt")
//	//if err != nil {
//	//	t.Errorf("Error not expected. %s", err.Error())
//	//}
//
//	expected := "Lorem ipsum"
//	if !strings.Contains(parsed, expected) {
//		t.Errorf("Expected string to be present:  %s", expected)
//	}
//
//	actualCharCount := len(parsed)
//	if actualCharCount > templating.MaxCharCount {
//		t.Errorf("String too long: %d is maximum. Actual: %d", templating.MaxCharCount, actualCharCount)
//	}
//}
//
//const longTranscript = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum rutrum neque felis, ut laoreet purus tempor vitae. Mauris gravida sapien in feugiat pharetra. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Nunc sagittis at dolor non venenatis. Etiam hendrerit ex vel ipsum euismod suscipit. Morbi non tellus sit amet sem suscipit malesuada sed a ex. Sed auctor feugiat tortor, nec fringilla urna. Integer pretium lacus scelerisque libero rhoncus, fermentum interdum augue vulputate. Maecenas vitae augue mauris.
//
//Suspendisse et erat et enim gravida fermentum. Quisque in purus non risus cursus mattis. Proin venenatis est eu tortor finibus, id gravida ex bibendum. Quisque rhoncus posuere sapien, nec placerat nibh vulputate non. Etiam sed velit a quam pretium fringilla sed nec massa. Suspendisse vel bibendum urna, non auctor tellus. Suspendisse in orci maximus, euismod tortor et, egestas erat. Nam maximus iaculis ex in venenatis. Suspendisse lobortis accumsan posuere.
//
//Sed suscipit scelerisque tellus. Quisque hendrerit libero eget enim gravida condimentum. Vivamus non erat purus. Pellentesque porta lacinia ante, sit amet porta leo sollicitudin vel. Aliquam semper a ante vel pellentesque. Fusce suscipit sollicitudin nulla, nec viverra nulla ornare quis. Fusce vel nisl rutrum, viverra nisi sed, sagittis metus. Morbi mauris lacus, posuere ac congue sed, commodo suscipit orci. Etiam pharetra ullamcorper mauris et consectetur. Aenean vulputate eros vitae dolor fermentum fermentum eget ut eros. Interdum et malesuada fames ac ante ipsum primis in faucibus. Pellentesque vel ornare nunc, sed placerat mi. Etiam pharetra, ante tristique aliquam pretium, purus mauris auctor elit, varius mollis leo elit id justo. Integer risus ex, mattis vitae nisl sit amet, tincidunt ornare lorem. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse non massa at quam egestas feugiat.
//
//Pellentesque et urna vel ipsum dapibus sodales. Nulla feugiat blandit odio eu pharetra. Nunc iaculis ligula tincidunt lobortis tempus. Donec semper luctus ornare. Ut dapibus bibendum dolor, eu luctus nibh gravida eget. Nullam fringilla magna lorem, ac tristique velit vehicula sagittis. Maecenas euismod felis et massa eleifend vulputate. Curabitur sit amet urna ut ex maximus rutrum vitae sed ipsum. Integer porttitor laoreet sem vitae ultrices. Donec pharetra feugiat lacinia. Suspendisse ac orci vitae felis tincidunt cursus. Fusce eu erat laoreet, volutpat odio sed, tempor odio. Nam felis neque, ornare id tellus vitae, ullamcorper facilisis purus. Nunc sapien leo, euismod eget augue ac, gravida suscipit quam. Nam nec urna erat. Aliquam sed mi id nisi tristique lobortis.
//
//Donec eget vestibulum tellus. Curabitur magna nibh, tincidunt sit amet mi eu, vulputate feugiat ipsum. Quisque elit quam, sagittis quis feugiat sed, posuere id magna. Proin id blandit risus. Nam eu eros ut sapien pulvinar ultricies. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Cras ac nunc justo. Vivamus id congue tortor. Pellentesque hendrerit est leo, sed lobortis turpis pellentesque vel. Mauris sodales ullamcorper lectus feugiat faucibus. Proin nec ante a risus tempus luctus.
//
//In ornare massa vel odio ornare dapibus. Ut nec blandit nunc. Mauris auctor luctus massa ultrices sodales. Vivamus volutpat arcu ut nulla sagittis, non malesuada lorem viverra. Curabitur ante lectus, iaculis eu dolor ac, interdum aliquam magna. Nulla ac faucibus ex, sed pharetra mi. Aliquam metus lectus, tempor quis tristique at, malesuada eu augue. Maecenas sit amet orci quis nisi sodales mollis eu eu ipsum.
//
//Donec faucibus tortor vel nunc blandit maximus. Aliquam ante ipsum, dapibus eget tortor sit amet, sollicitudin blandit erat. Donec ut neque nec quam congue bibendum. Aliquam lacinia id lacus eget gravida. Ut convallis maximus efficitur. Sed ultricies erat quis tortor faucibus fermentum. Fusce magna nunc, vestibulum sodales urna suscipit, gravida tempor nisl. Praesent nec convallis nisi, in aliquet nulla. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.
//
//Etiam vulputate nulla eget nisl molestie luctus id ac lectus. Donec nec felis ipsum. Ut tincidunt lectus diam, a posuere justo gravida sed. Integer sed massa quis enim volutpat convallis. Aliquam vestibulum nisl a bibendum laoreet. Phasellus mi risus, convallis ac enim sit amet, ultrices elementum ligula. Integer viverra mattis tincidunt. Proin pellentesque ornare lacus vitae mattis. Aenean a turpis tortor. Duis iaculis diam nec consequat efficitur. Maecenas vel porta lacus. Vivamus ullamcorper ac neque at luctus. Duis accumsan tellus vitae odio cursus, pretium sollicitudin enim fermentum. Cras mollis nisl ac lectus pellentesque accumsan.
//
//Nulla vitae risus varius, malesuada dolor at, tincidunt magna. In nibh tortor, elementum et vehicula eget, placerat at velit. Praesent mi felis, ullamcorper in nulla vel, tincidunt luctus lectus. Cras pellentesque porta dolor, a dapibus quam vulputate accumsan. Donec laoreet, ex eu pellentesque pharetra, diam purus condimentum sem, ut porttitor sem nunc vel odio. Maecenas condimentum nisl odio, ac tempus ex congue quis. Nullam volutpat ac augue eget euismod. Vivamus commodo convallis ligula, quis rhoncus magna vestibulum ut. Morbi a felis ultrices leo tincidunt condimentum.
//
//Pellentesque condimentum lobortis enim. In metus justo, commodo ut mi vitae, elementum porta nisl. Cras ligula ipsum metus.`
