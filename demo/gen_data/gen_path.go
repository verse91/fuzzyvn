//go:build ignore

package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

const (
	OutputPath = "test_paths_100k.txt"
	MaxPaths   = 100000
	MinWin     = 6
	MaxWin     = 18
)

var (
	wordRe = regexp.MustCompile(
		`[a-zA-Z0-9àáạảãâầấậẩẫăằắặẳẵ` +
			`èéẹẻẽêềếệểễ` +
			`ìíịỉĩ` +
			`òóọỏõôồốộổỗơờớợởỡ` +
			`ùúụủũưừứựửữ` +
			`ỳýỵỷỹđ]+`)
	spaceRe = regexp.MustCompile(`\s+`)
)

func alphaBucket(word string) string {
	if word == "" {
		return "khac"
	}

	r := []rune(word)[0]
	switch {
	case r >= 'a' && r <= 'f':
		return "người"
	case r >= 'g' && r <= 'l':
		return "tai nạn"
	case r >= 'm' && r <= 'r':
		return "cầu thủ"
	case r >= 's' && r <= 'z':
		return "siêu việt"
	default:
		return "khac"
	}
}

func shortHash(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])[:6]
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = norm.NFC.String(s)
	return s
}

func main() {
	rawText := `
	 Bên trong xe là tiếng nhạc ầm ĩ, mùi rượu cay nồng và 4-5 người đàn ông sẵn sàng thách thức bất kỳ ai dám đối đầu Nạn nhân may mắn không bị thương, những người trên xe sau đó bỏ đi với vẻ bất cần, còn anh Thanh ở lại hỏi han người bị nạn và cho số điện thoại để sẵn sàng ra làm chứng khi cần.
Hoàng Thanh chia sẻ, việc anh làm bị nhiều người cho là liều lĩnh, nhưng bản thân anh lúc đó cũng không đủ thời gian để suy nghĩ, cân nhắc nhiều mà chỉ nghĩ rằng cần phải làm gì đó giúp đỡ người bị nạn và yêu cầu người gây ra lỗi có trách nhiệm với hành vi của mình.
Anh Nguyễn Hoàng Thanh (trái) bắt tay phóng viên Báo An ninh Thủ đô Chàng sinh viên làm việc thiện ngại gặp vì quần áo bẩn Đầu tháng 1-2018, một vụ tai nạn giao thông khá nghiêm trọng xảy ra khiến nạn nhân là một người đàn ông trung niên nhập viện và cần phải phẫu thuật.
Tôi nghĩ, cuộc sống này khó tránh khỏi những nốt trầm, những màu buồn u tối, những khoảnh khắc giận dữ, bức xúc Nhưng bỏ qua mọi thứ tiêu cực đó, vẫn có những con người tràn đầy sức sống, sẵn sàng lan tỏa năng lượng tích cực thông qua nụ cười, ánh mắt và hành động ý nghĩa của họ.
Mình thấy có rất nhiều bạn bị mụn chỉ quan tâm đến dập lửa hoặc có 1 số thì tìm kiếm đến các loại thực phẩm chức năng được quảng cáo là detox thải lọc độc tố giải trừ mụn... nhưng tuyệt nhiên các cách rẻ tiền tự thân làm được thì không mấy ai chú ý tới và có vẻ thờ ơ coi thường không mấy thực hiện.
Từng bị trầm cảm, xa lánh mọi người vì mụn mọc chi chít Thời gian trước cách đây 2 năm mình rơi vô trầm cảm, xa lánh tất cả mọi người lui về ở ẩn nơi hẻo lánh không giao du, buông xuôi tất cả và hoàn toàn bế tắc khi nghĩ về tương lai vì mụn mọc nhiều trên mặt.
Mỗi ngày mình đi làm và vận động cơ thể bằng cách đi bộ rất nhiều, cơ thể mình toát mồ hôi nhễ nhại, rất rất nhiều mồ hôi trên toàn thân, từ cả bẹn, 2 bên đùi, nách, lưng, ngực mồ hôi như tắm ướt đẫm dù mình không hoạt động nặng nhọc gì chỉ đi bộ nhiệt tình thôi.
Hiện nay đời sống khá tất bật nhưng hầu như mọi người đều bận rộn theo kiểu ngồi văn phòng có điều hòa lạnh rồi đi về nhà tối tối chăm chỉ thoa bôi trát kem dưỡng và rất ít chị em để ý đến khâu làm sao để sản xuất được nhiều mồ hôi nhất mỗi ngày.
Họ cứ thế mà hy sinh mọi thứ, miễn sao là người họ yêu được vui vẻ và hạnh phúc, không cần hồi đáp, thậm chí chấp nhận việc bản thân bị đối xử tệ bạc không ra gì, chấp nhận trở thành một người thừa trong một cuộc tình bất kỳ nào đó.
Vậy thì đàn bà à, hãy nghe tôi nói một câu, sinh ra đã phận thiệt thòi, đàn bà đừng hy sinh cho những điều vô nghĩa và đừng đặt niềm tin của mình vào một người đàn ông mà chúng ta không chắc chắn rằng họ sẽ đem lại cho ta được hạnh phúc như bản thân chúng ta kỳ vọng.
Để bên cạnh người đàn ông xứng đáng, đàn bà hy sinh thời gian, tự do, nhan sắc, hy sinh cả những dự định, vì gặp người đàn ông ấy mà đột ngột rẽ ngang, người ngoài nhìn vào cho rằng đàn bà mù quáng, ở trong thâm tâm thì đàn bà thấy mình mãn nguyện.
Phụ nữ biến thành đàn ông Những bức ảnh về các trinh nữ thề sống như đàn ông trong bài được chụp bởi nhiếp ảnh gia người Mỹ Jill Peters Theo trang Amusing Planet, những người phụ nữ sống cuộc đời của một người đàn ông ở những vùng núi hẻo lánh thuộc miền Bắc Albania không chỉ phải từ bỏ hình dáng nữ tính của mình mà thậm chí còn phải bỏ tên cũ và thay vào đó là một cái tên của nam giới.
Antonia Young, nhà nhân chủng học, đồng thời là tác giả của cuốn Women Who Become Men: Albanian Sworn Virgins (tạm dịch: Phụ nữ biến thành đàn ông: Những trinh nữ gắn với lời thề ở Albania) giải thích thêm rằng, khi đã được gia đình định hướng, các bé gái sẽ được nuôi dạy ứng xử, ăn mặc như con trai, tổ chức, quán xuyến toàn bộ công việc trong gia đình... Những người đã thề sẽ thành con trai không bao giờ được trở lại thân phận phụ nữ và nếu vi phạm lời thề, họ sẽ bị trừng phạt bằng cái chết.
Họ không được phép đi bỏ phiếu, không được kế thừa tài sản, mua đất, kinh doanh, kiếm tiền, hoặc thậm chí hút thuốc, đeo đồng hồ... Những luật lệ này được gọi là Kanun thịnh hành suốt 5 thế kỷ từ thế kỷ 15 đến thế kỷ 20 và thậm chí đến nay vẫn chi phối trong khu vực.
Mặc dù không thể trải nghiệm niềm vui của một người phụ nữ chẳng hạn như lấy chồng, mang thai và sinh con đẻ cái thì sự tự do và những đặc quyền của đàn ông cũng đủ quyến rũ để nhiều phụ nữ tự nguyện sắm vai nam trọn đời.
Bé Trung sau ca phẫu thuật đã dần hồi phục Còn nước, còn tát Trước cánh cửa phòng cấp cứu Bệnh viện Đa khoa tỉnh Quảng Ninh, chị Cầm Thị Linh (trú tại xã Dực Yên, huyện Đầm Hà, tỉnh Quảng Ninh) không lúc nào ngồi yên một chỗ, chị đi đi, đi lại rồi khóc thầm cầu mong đứa con trai 7 tuổi Trần Ngọc Trung của mình bị TNGT được cứu sống.
Chị vẫn ngồi đấy để trực chờ y tá đưa thức ăn đến để phụ giúp tiêm vào ống dẫn cho con, miệt mài nói chuyện vì con không thể ngủ và luôn ngồi bóp nhẹ đôi bàn tay nhỏ bé để giữ chặt Trung ở lại sau vụ TNGT kinh hoàng.
Các bác sỹ đã khẩn trương sơ cứu chống sốc và tiến hành hội chẩn nhanh giữa các chuyên khoa xác định cháu Trung bị sốc đa chấn thương, chấn thương sọ não, tụ máu trong não, chấn thương bụng kín, vỡ lách độ IV-V, gãy xương đùi và vỡ xương chậu.
Bệnh nhi mới 7 tuổi lại bị đa chấn thương phức tạp, trường hợp này nếu đồng thời thực hiện phẫu thuật cắt lá lách và mổ sọ não thì nguy cơ dẫn đến tình trạng rối loạn đông máu là rất cao, đe dọa đến tính mạng ngay trong khi mổ cấp cứu.
Quyết định hiến xác từ khi còn sống Vào một ngày trời âm u giá rét, cơn mưa phùn lại càng khiến tiết trờ trở nên se lạnh hơn, chúng tôi đã tìm về ngôi nhà của chị Nguyễn Thị Phòng ở thôn Đồng Trạng (xã Cổ Đông, thị xã Sơn Tây, Hà Nội) gặp người phụ nữ hiến xác chồng cho y học và hiến những bộ phận trên cơ thể còn dùng được để cứu sống người khác.
Lúc đầu nghe xong, chị phản đối gay gắt nhưng sau được chồng tâm sự, giải thích chị hiểu ra được vấn đề nên đã đồng ý. Thẻ đăng ký hiến mô, tạng của chồng chị Phòng Lúc chưa bệnh tật gì anh ấy bảo, tớ không biết sống đến bao giờ nhưng khi nào tớ mất tớ cứ hiến tạng họ lấy được gì cứu ai thì càng tốt tốt.
Chị Phòng bảo, khi chồng chị xuống làm lễ tri ân về có nói: Họ làm lễ chu đáo lắm, sau này họ thiêu luôn cho sạch sẽ, mẹ mày cứ để anh ra đi như vậy không phải lăn tăn đâu, ở đấy còn sướng hơn ở nhà, lại mát mẻ nữa.
Thế nhưng chị Phòng đã rất kiên quyết và từ tốn giải thích đó là nguyện vọng của anh Trường muốn cái chết của mình có ý nghĩa mọi người đều đồng ý Vậy là thi thể của anh Trường được chuyển về Học viện Quân y. Khi được hỏi về phần thân thể còn lại của anh Trường bao giờ sẽ trở về với gia đình chị Phòng chia sẻ, chị nhất trí để xác chồng lại cho 5 năm sau đó sẽ cho chồng quay trở về nhà.
Trên truyền hình có trường hợp, sau khi con trai bị chết não do tai nạn giao thông, người mẹ ở Quốc Oai đã hiến một số bộ phận trên cơ thể con cho y học và cứu được 5 người, trong đó tim của họ cứu được một chiến sĩ ở hải đảo.
Thế nhưng, chị Phòng đã giải thích và thuyết phục mọi người, đó là nguyện vọng của anh muốn cái chết của mình có ý nghĩa Vậy là hơn 10h đêm ngày 3/10/2017, thi thể của anh Trường được chuyển về học viện Quân y. Tôi cũng bảo với các con, sau này mẹ cũng sẽ nối bước của bố, sẽ hiến mô tạng cho y học , chị Phòng khẽ ngước nhìn di ảnh chồng rồi mỉm cười.
Về tâm linh thì tiền duyên là mối duyên tình mang yếu tố say đắm, đam mê hay hận thù trong vòng nhân quả từ kiếp trước chưa được giải quyết giữa những người đã từng luyến ái với nhau, vấn đề này không thể qua loa kiểm chứng bằng ngoại cảm đồn đại, nở rộ trước giờ.
Hiện nay, những người bị phán có tiền duyên đều do thầy đồng, bà cốt lợi dụng tâm lý của những cô gái, chàng trai đang sốt ruột chuyện lập gia đình để bịa ra chuyện có tiền duyên oan trái hay âm hồn oan tình báo oán nhằm mục đích trục lợi, kiếm sống bản thân.
Lý giải về tiền duyên và việc cắt tiền duyên , TS.Vũ Thế Khanh, Tổng Giám đốc Liên hiệp Khoa học Công nghệ Tin học Ứng dụng UIA, người có kinh nghiệm 20 năm trong lĩnh vực khoa học góc nhìn tâm linh, cho rằng: Tiền duyên chính là luật nhân quả của mỗi người.
Theo Ths.Vũ Đức Huynh, tác giả cuốn sách Con người với tâm linh , thất tình, lục dục là bảy thứ tình cảm và sáu thứ dục vọng ở con người, trong đó có ba thứ dục vọng là nhu cầu sống của bản thân là: Ăn uống, tình dục và ngủ nghỉ.
Các bài thuốc trị rối loạn tiêu hóa ở trẻ: cháo rau sam; cà rốt; nước nụ vối;.... Nếu trẻ đi cầu nhiều lần, phân sống, mùi chua, biểu hiện suy dinh dưỡng nên kiện tỳ tiêu thực cho trẻ: Đảng sâm, hoài sơn, ý dĩ mỗi thứ 6g, nhục đậu khấu, trần bì, mạch nha, hậu phác mỗi vị 4g, sắc ngày 1 thang chia nhiều lần cho trẻ uống trong ngày.
Giảm 3kg trong vòng 10 ngày với chế độ ăn kiêng dành riêng cho người Châu Á Tiến sĩ Colin Campbell đến từ Trung tâm nghiên cứu Trung Quốc đã nghiên cứu và tìm ra bí quyết giảm cân của người châu Á. Theo ông, nếu áp dụng chế độ ăn sau đây, bạn hoàn toàn có thể giảm được 3kg trong vòng 10 ngày.
Chuối được lựa chọn bởi loại quả phổ biến này chứa một loại chất xơ đặc biệt có tác dụng kháng tinh bột có trong chuối sẽ bị lên men trong dạ dày tạo ra sản phẩm giúp đốt cháy 20-25% chất béo, do đó làm giảm năng lượng nhập vào cơ thể, tránh sự tích tụ mỡ thừa.
Năm ngoái, con gái của một cựu binh Mỹ từng tham chiến ở Việt Nam và chịu hậu quả nặng nề của Hội chứng chiến tranh Việt Nam, cô ấy muốn sang Việt Nam, tìm cách giúp đỡ các nạn nhân chất độc da cam để linh hồn người cha quá cố của mình được thanh thản.
Và, khi biết được đầu đuôi câu chuyện, rằng anh ta, khi thấy một người chia sẻ trên mạng rằng sẽ xăm hình lá cờ Việt Nam lên ngực nếu đội tuyển U23 Việt Nam giành được chức vô địch bóng đá châu Á, Daniel Hauer đã mỉa mai bằng một câu nói khiếm nhã, bất kính tới Đại tướng Võ Nguyên Giáp; chính tôi cũng cảm thấy nghẹt thở vì tức giận.
Anh Võ Trung, người cháu nội của cố Đại tướng Võ Nguyên Giáp đã phản ứng trên Facebook: Thường chuyện gì tôi cũng có thể bỏ qua nhưng việc lần này thì không, tôi cũng quen rất nhiều bạn bè nước ngoài và văn hóa khác biệt có những trò đùa như thế nào tôi cũng biết, nhưng người này thì khác hoàn toàn bởi chúng ta phần lớn đều đã xem qua những clip được chia sẻ nhiều về văn hóa Việt Nam và anh ta cũng là một người rất hiểu về văn hóa nơi đây.
Có lẽ chính bởi một gia đình được xây nên bởi nền tảng là tình cảm hạnh phúc, sự kính yêu và trân trọng, yêu thương lẫn nhau như vậy nên không có gì lạ khi đại gia đình thành công, với những người con thành đạt, những người cháu học hành giỏi giang và hết mực hiếu thảo.
Chỉ là một đứa trẻ bán báo dạo, kiếm từng đồng bạc lẽ để lo cho cuộc sống, thậm chí kết nghĩa anh em với nhiều giang hồ máu mặt thời đó như Đại-Tỳ-Cái-Thế nhưng sư Thiện vẫn không bị ảnh hưởng tiếng xấu để đời mà vẫn giữ được bản tính hiền lành và luôn ghi nhớ lời mẹ dặn Giấy rách phải giữ lấy lề Từ chối lời đề nghị của các đại ca Như đề cập ở kỳ trước, sư cô Diệu Thiện tên thật là Nguyễn Thị Sự cùng anh trai tên Việt phiêu bạt khắp nơi để kiếm sống sau khi mồ côi cả cha lẫn mẹ.
Sư Diệu Thiện lúc còn trẻ Lúc bấy giờ sư cô gặp biết bao nhiêu cám dỗ của cuộc đời khi theo chân những đại ca khét tiếng một thời như Trần Đại (Đại Cathay), Bảy Si, Lâm Chín Ngón, Năm Cam.Thế nhưng, theo sư cô Thiện: Trong giới giang hồ lúc đó cái tên uy lực nhất trong Tứ đại thiên vương là Tín Mã Nàm, một tên giang hồ gốc Hoa quản lý khu vự Chợ Lớn (quận 5).
Trao đổi với PV Chất lượng Việt Nam Online , một chuyên gia điện lạnh làm việc tại siêu thị điện máy Trần Anh cho hay, những nguyên nhân chính dẫn đến điều hòa bị hỏng là do các gia đình sống tại chung cư có ban công nhỏ, hẹp, khi sử dụng điều hòa lại chọn sai chế độ/chức năng trên điều khiển.
Lắp mặt lạnh quá cao hoặc không đúng vị trí sử dụng của người dùng dẫn đến tình trạng nơi cần mát không mát, nơi không cần mát lại mát Ngoài ra, chuyên gia cũng khẳng định, việc bảo dưỡng điều hòa trước mùa nóng hàng năm vô cùng cần thiết.
Đối với nhà có trẻ nhỏ nên duy trì nhiệt độ lý tưởng từ 28 29 độ C. Đối với nhà không có trẻ nhỏ thì nên duy trì nhiệt độ từ 22 23 độ C vào ban ngày, ban đêm nên duy trì nhiệt độ từ 25 26 độ C hoặc tắt hẳn.
Việc niệm Phật với tâm từ bi hỷ xả là hoàn toàn tốt, cho dù ở bất kỳ hoàn cảnh nào, nó gieo trong ta một niềm an lạc, tình thương yêu từ đó phát sinh hạt bồ đề chăm làm thiện, lánh việc ác dẫn đến cơ thể luôn cân bằng để vững bước trên đường đời.
Biết bao nhiêu người hôm nay đã không thấy được vầng thái dương của ngày mai, biết bao nhiêu người hôm nay đã trở thành tàn phế, biết bao nhiêu người hôm nay đã đánh mất tự do, biết bao nhiêu người hôm nay đã trở thành nước mất nhà tan.
Phong thủy ứng dụng và tâm linh đều cho rằng, những người mà kiếp này may mắn có địa vị cao quý, được xã hội trọng vọng làm quốc vương hay chức cao đại thần, là người có quyền, có thế thì kiếp trước đều là những người lễ phép, biết kính trọng Phật, Pháp, Tăng mà đến.
Người nào kiếp này có cá tính điềm đạm, cư xử bình tĩnh, hành xử không bao giờ hấp tấp vội vàng, cả trong nói năng và trong hành động đều rất cẩn trọng, biết chừng mực thì ắt hẳn kiếp trước đều là những người đã từng tu thiền định, tâm tưởng thanh tịnh.
Người kiếp này tài năng và thông suốt Pháp, thậm chí có thể thuyết giảng, đồng thời hóa độ người u mê hay ngốc nghếch và hiểu được, biết trân quý lời nói và tự động truyền rộng Phật pháp ra ngoài để người người trong chúng sinh cùng thấu hiểu.
Phẫu thuật thẩm mỹ không còn xa lạ gì với phái đẹp, có muôn vàn lý do chị, em tìm đến sự can thiệp của dao kéo trong việc thay đổi diện mạo của mình, người thiếu tự tin về nhan sắc thì tìm đến sự can thiệp của dao kéo để mong muốn được ưa nhìn hơn, người đẹp rồi lại muốn được đẹp hoàn hảo hơn nữa.
Lời cảnh tỉnh cho phái đẹp Những lời có cánh được quảng cáo trên web của TMV với đội ngũ bác sĩ uy tín được cấp chứng chỉ tại Hàn Quốc, hệ thống TMV có thiết bị và công nghệ Hàn Quốc, cam kết bảo hành100% kết quả phẫu thuật thẩm mỹ cho khách hàng... Những lời quảng cáo "ngọt như rót mật vào tai đã câu kéo các chị em đổ xô đi làm đẹp.
Hành trình đau đớn biến thành Thiên nga của khách hàng tại TMV Trao đổi với phóng viên, một khách hàng của một TMV - chị L - ngậm ngùi trong nước mắt với "khuôn mặt méo mó, dở khóc dở cười : "Sau khi được tư vấn tôi đã tiến hành phẫu thuật thẩm mỹ tại một TMV tại Thành phố Hồ Chí Minh.
Khi Vịt Không Thành Thiên Nga' TMV lại đổ lỗi do cơ địa của khách hàng Hình ảnh trước và sau khi phẫu thuật của chị L Theo chị L tường trình: sau khi phẫu thuật, mặt chị ngứa và bị sưng tê cứng hết đầu, còn mắt thì không nhắm được.
Một tháng phẫu thuật và gây mê 3 lần vẫn không thành thiên nga Chị L bức xúc : "Biết tình hình của tôi không ổn TMV cố tình bưng bít không cho tôi biết sự thật mà toàn nói quanh co, và đưa ra một số dẫn chứng bằng hình ảnh của các khách hàng trước đó bảo tôi cứ yên tâm và bình tĩnh đây là các trường hợp cũng giống như tình trạng của tôi, nhưng sau 3 tháng sẽ khỏi hoàn toàn nếu phẫu thuật lại.
Có lẽ nói về cách dùng đũa thì có lẽ không ai hơn người Á châu, tuy nhiên có lẽ nhiều người Việt vẫn có những thói quen dùng đũa không tốt cho tài vận, thậm chí là mang điềm gở cho chính người dùng mà ta không hay để ý.
Nối đũa nhau Khi gắp thức ăn cho người khác, ngoài việc phải trở đầu đũa để giữ vệ sinh cho người nhận, bạn còn phải để ý gắp thức ăn bỏ hẳn vào chén của họ và tránh việc nối đũa , tức là chuyền thức ăn từ đũa mình sang đũa người khác.
Đời của bạn là do bạn quyết định Phật dạy , ba mẹ sinh ra bạn nhưng không thể nào ở bên bạn suốt đời, chồng bạn ở bên bạn hôm nay nhưng cũng chưa chắc còn bên bạn ngày mai, con cái bạn vẫn còn nhỏ trong lòng bàn tay nhưng tương lai cũng sẽ đi bên tay người khác.
Chạy đua với thời gian Tiền tài và danh vọng là không chờ đợi Theo lời Phật dạy, nhân lúc đang còn trẻ, dũng cảm bước đi, nghênh đón phong sương gió mưa, tôi luyện bản thân, có thể độ lượng, có thể nhìn xa trông rộng, thì hạnh phúc mới đến.
Nếu bạn chờ người khác làm vỡ bạn từ bên ngoài, thì nhất định bạn sẽ là món ăn của người khác; nếu bạn có thể đánh vỡ chính mình từ bên trong, như vậy bạn sẽ thấy rằng mình đã thực sự trưởng thành, cũng giống như là được tái sinh.
Biết bố thí, cúng dường Theo luật nhân quả, sự trộm cắp, bủn xỉn, ích kỷ, tham lam là nhân đưa đến sự nghèo túng thì ngược lại, để thoát khỏi điều này chúng ta hãy biết bố thí, cúng dường, san sẻ vật chất cho những người nghèo khổ hoặc hộ trì chánh pháp.
Người phụ nữ mạnh mẽ thông minh chắc chắn không bao giờ làm điều này khi yêuNgười phụ nữ mạnh mẽ thông minh chắc chắn không bao giờ làm điều này khi yêu Đây chính là điều mà người phụ nữ thông minh sẽ làm sau khi ly hônĐây chính là điều mà người phụ nữ thông minh sẽ làm sau khi ly hôn Thói quen trong hành xử không ngờ khiến bạn bị người khác xem thường quá nhiều người chẳng hayThói quen trong hành xử không ngờ khiến bạn bị người khác xem thường quá nhiều người chẳng hay Theo An Nhiên/Khoevadep.
Con người bị chi phối bởi những chu kỳ sinh học, chẳng hạn như chu kỳ trí tuệ, sức khỏe và tâm lý Khoa học hiện đại cũng chỉ ra con người bị chi phối bởi những chu kỳ sinh học, chẳng hạn như chu kỳ trí tuệ, sức khỏe và tâm lý.
Từ lúc mang con trong lòng, sinh con ra cho đến lúc con trưởng thành, cha mẹ phải vất vả nhọc nhằn, tốn hao biết bao mồ hôi nước mắt, công sức, dành hết tâm tư tình cảm cho con, chỉ mong con nên vóc nên hình, khôn lớn thành người.
Vì việc làm có vẻ ngớ ngẩn nên mọi người mặc kệ cậu bé với hành động kì quái của mình, cho đến ngày đống đá nhô lên thành ụn đất cao, người dân trong làng mới cùng nhau góp công, góp của xây cầu để việc đi lại thuận lợi và an toàn hơn.
Truyền thuyết kể lại rằng, nhờ chiếc gối Âm Dương, Bao Công gặp được một vị trời hỏi về sự việc kì lạ này, được trả lời rằng: Trong quá khứ, cậu bé kia từng là kẻ đại gian ác, giết người cướp của, hiếp đáp kẻ yếu, nên phải trả nghiệp bằng ba kiếp què, mù, và bị sét đánh chết.
Và chắc chắn nghiệp tốt sẽ không mất đi, như Kinh Nhân Quả có nói: Người nay giàu có, vì đời trước thường hành Bố Thí Người có tướng mạo xinh đẹp, vì đời trước có tâm cung kính Do vậy, khi thấy một người sống tốt, luôn hành động vì người khác mà vẫn gặp khổ, hãy hiểu rằng đó là sự hiển bày của nhân quả mà thôi.
Nếu có thể tóm tắt nhất lời phật dạy chỉ trong 4 câu, thì không gì hơn bài kệ thứ 183 trong Kinh Pháp Cú: Không làm các điều ác Làm tất cả điều lành Tịnh hóa tâm ý mình Là lời Chư Phật dậy (Chư ác mạc tác - Chúng thiện phụng hành - Tự tịnh kỳ ý - Thị chư Phật giáo) Người tốt khổ vì trong lòng còn có ác tâm Chữ Khổ trong đạo Phật không chỉ là cảm giác khổ sở tại thân tâm, mà còn bao gồm cả những giao động rất nhẹ của sự Không Thỏa Mãn, Không Yên Ổn, hay còn gọi là Bất Toại Nguyện ở trong lòng.
Một người cảm thấy mình rất tu tâm tích đức, thường xuyên làm điều tốt, chả hại ai bao giờ, ấy vậy mà cuộc đời mãi cứ long đong vất vả, không được bằng bạn bằng bè, một lần lên chùa vãn cảnh, quá buồn cho phận đời mình liền tiến đến hỏi một vị sư: Thưa thầy, vì sao con sống tốt, sống thiện, mà đời con cứ khổ mãi chưa thấy khá lên?
Từ khi rừng Pác Ngòi bị chặt phá, cây to bị xẻ thành ván đóng thành tủ, ghế, cây nhỏ thành củi, thành than bán cho những thợ rèn, thảm thực vật còn lại được người làng gom đốt, mỗi nhà chia nhau một khoảnh nương trồng lúa nếp nương, sắn, nước cũng theo cây lặn sâu vào lòng đất.
Mấy lần họp bản bàn đến việc lấy nước, nhiều người có ý kiến nên lấy nước từ lũng Rưa Bắc về hợp với mỏ nước uống chảy từ trong núi và suối Pác Ngòi ở đầu bản, nước từ đây sẽ tỏa đi nhiều hướng, không phải chờ nước trời.
Anh đã thưa hỏi nhiều vị đạo sư nổi tiếng mà chưa một vị nào giải đáp được thỏa mãn tâm tư, nguyện vọng của anh, nghe đồn Phật là bậc xuất trần thượng sĩ, có thể giải quyết được mối nghi ngờ của nhiều người bất kể là tín đồ tôn giáo nào, anh đã tìm đến đức Phật.
Sau khi đảnh lễ và vấn an đức Thế Tôn xong, anh cung kính ngồi sang một bên thưa hỏi đức Phật: Vì cớ sao có sự bất công và sai biệt quá lớn của tất cả chúng sinh trên thế gian này, kẻ quý phái cao thượng, người hạ liệt thấp kém, người sống lâu, kẻ chết yểu, người giàu sang, kẻ nghèo khổ, người nhiều bệnh, kẻ ít bệnh, người quyền cao chức trọng, kẻ nô lệ thấp kém, người đẹp đẽ dễ thương, kẻ xấu xí khó nhìn, người thông minh sáng suốt, kẻ ngu dốt tối tăm?
Đức Phật trả lời quá xúc tích và cô đọng làm chàng thanh niên không thể hiểu rõ nghĩa lý sâu xa, nên mới yêu cầu đức Phật giải thích cụ thể từng chi tiết : Phật dạy: Tất cả mọi sự sai biệt giữa con người và con người là do nghiệp của ta đã tạo ra từ thân miệng ý, tâm suy nghĩ chân chính, miệng nói lời thiện lành, thân đóng góp sẻ chia, thì được hưởng quả an vui hạnh phúc; ngược lại, gieo nhân xấu ác thì bị quả sa đọa khổ đau, không ai có quyền xen vô chỗ này để định đoạt và sắp đặt, nên có người tốt, kẻ xấu là do mình tạo ra.
Trong trường hợp phúc báo nên được hưởng sắp hết mà việc ác đã làm lại tích tụ quá nhiều, lúc này là thời cơ ác báo đã chín muồi , phúc báo kết thúc và người này bắt đầu phải chịu ác báo như sống thê thảm, đột nhiên bị bệnh tật, tai nạn bất ngờ mà chết, hối hận không kịp.
Sắp tới, các bác sĩ sẽ mổ lấy chất làm đầy mũi, sau đó tạo hình mũi cho chị K. Tương tự, vì muốn có một khuôn mặt đầy đặn, chị V.T.H (29 tuổi, ngụ quận 10, TP HCM) lên mạng tìm kiếm một cơ sở thẩm mỹ để thực hiện phẫu thuật.
Một phụ nữ bị biến chứng vùng cổ khi đi thẩm mỹ Đừng biến mình thành "vật thí nghiệm" Trao đổi với Báo Người Lao Động, PGS-TS Đỗ Quang Hùng - Tổng Thư ký Hội Phẫu thuật thẩm mỹ TP HCM, Trưởng Khoa Tạo hình thẩm mỹ BV Chợ Rẫy - nhấn mạnh: "Để có thể hành nghề phẫu thuật tạo hình thẩm mỹ thì sau khi học 6 năm ở trường y, bác sĩ được đào tạo thêm 3 năm về chuyên ngành thẩm mỹ, sau đó thực hành ít nhất 2 năm.
Kim Chi cho rằng mình là người dám yêu, dám bỏ nên chọn được nhà chồng tốt Tuy nhiên, lại có luồng ý kiến cho rằng, cô vô tư đến mức vô ý. Dù nhà chồng thoải mái đến đâu, con dâu cũng nên hiểu trách nhiệm của mình mà làm cho đúng.
C kể, tối qua 2 vợ chồng đi siêu thị, chồng C đã loanh quanh chỗ quầy bầy rượu để chọn quà Tết cho ông ngoại: Anh đứng 1 lúc thì kêu, em xem có chai nào đẹp đẹp nhìn sang xịn 1 tí để sang tháng mua về ông ngoại biếu tết.
Tôi yêu cô ấy và đau xót khi nhìn thấy cô ấy bị người ta ghẻ lạnh, muốn cô ấy bỏ chồng, đến với tôi, chúng tôi sẽ làm lại từ đầu với nhau, tôi sẽ yêu thương cô ấy, yêu thương con cô ấy như chính con của mình, nhưng cô ấy từ chối.
Thực ra, gia đình Welde đã có giấy phép chăm sóc gấu gồm gấu đen, gấu Bắc cực... trong 91 năm qua, nhưng truyền thống gia đình bị lung lay khi vào năm 2016, chồng của bà Monica, ông Johnny III, đột ngột qua đời sau một cơn đau tim khi ông vừa 60 tuổi.
Cô gái xinh đẹp Manpreet - người trưởng tộc nắm giữ quyền thừa kế Trưởng tộc Mannu với vẻ ngoài điển trai Sinh trưởng trong gia đình có bà nội là người mang nặng tư tưởng trọng nam khinh nữ, nên ngay khi lọt lòng, Mannu (tên thuở nhỏ của Manpreet) của Bí mật người thừa kế đã được mẹ cải trang, nuôi dưỡng và mặc định là cậu-con-trai-đích-thực.
Trước sự quản thúc, giáo dục nghiêm khắc từ mẹ, cùng sự đanh ác của bà nội và các thành viên trong gia tộc, Mannu hiểu được nếu cậu không phải là con trai thì gia đình - trong đó có mẹ, hai chị gái và cả cậu , sẽ phải chết.
Với cái tên Mapreet, Mannu đã được trở lại với chính mình với những rung động đầu đời 17 năm trôi qua, việc bảo vệ bí mật về cậu-con-trai và cũng là người thừa kế sáng giá của gia tộc là hành trình vô cùng gian khó nếu như không muốn nói đầy máu và nước mắt với những người lỡ nắm giữ bí mật này.
Và ở tuổi đôi mươi với những đổi thay hình thể, Mannu nhận ra cậu thật sự khát khao được làm con gái để diện lên người những bộ cánh yêu thích và khoe những đường cong hình thể đáng mơ ước dù chỉ 1 ngày như lời hứa của mẹ từ 10 năm trước.
Đôi bạn thân thơ ấu - Mannu và Raj hội ngộ sau 10 năm Và cuối cùng, với cái tên Manpreet, Mannu đã được trở lại với chính mình - một cô gái xinh đẹp, quyến rũ và quyết định này đã giúp cô gặp lại Raj Bajwa - cậu bạn thân năm xưa, kéo theo đó là những rung động của tình yêu đầu đời... Đây cũng là lúc Mannu phải đấu tranh giữa trái tim và lý trí - trở lại với chính mình để giữ lấy tình yêu hay vẫn sẽ là một người thừa kế sáng giá của gia tộc?
Soumya nhân hậu - nạn nhân của sự kỳ thị giới tính Và tất cả những gì người khác cảm nhận được ở Soumya - cô gái dị giới trong bộ phim "Hai số phận" chính là một cô gái kiên cường đầy quyến rũ Sinh ra đã mang giới tính không trọn vẹn và chịu sự kỳ thị của xã hội, thậm chí bị ngay chính cha mình tìm cách chôn sống, nhưng nghị lực đã giúp Soumya chiến thắng tất cả.
Dù luôn được chồng yêu thương, chiều chuộng nhưng Soumya hiểu được, để bảo vệ hạnh phúc của bản thân thì không gì khác ngoài việc phải chấp nhận quyết định của nhà chồng bằng việc để Harman - người chồng đầu ắp tay gối của mình, kết hôn với chính em gái mình.
Suốt từ tối qua, dư âm chiến thắng lịch sử của đội tuyển U23 Việt Nam tại giải U23 Châu Á đang là đề tài được tất cả mọi người bàn tán thì 1 nàng dâu trẻ đã tâm sự trong những phút hân hoan chiến thắng như vậy nhưng cô vẫn thấy khó chịu vì có người bố vừa bảo thủ vừa cổ hủ ghê gớm.
Suốt từ tối qua, dư âm chiến thắng lịch sử của đội tuyển U23 Việt Nam tại giải U23 Châu Á đang là đề tài được tất cả mọi người bàn tán Không cãi được là ông chửi, ông nói đến bao giờ người ta chịu im hoặc nhận sai cho ông đỡ nói nhưng ông vẫn nói, nói mãi.
TAND quận 6, TP.HCM cho biết cơ quan này đã thụ lý đơn kiện vụ đòi bồi thường tổn thất sức khỏe giữa nguyên đơn là chị Nguyễn Thị Loan (sinh năm 1994, quê Đắk Lắk), bị đơn là thẩm mỹ viện Hà Anh ở đường Đặng Nguyên Cẩn, phường 13, quận 6 do bà Tô Thị Ngân Hà làm đại diện.
Tiêm đầy mũi gây hư mắt Theo đơn của chị Loan, từ tháng 12/2016, chị đã đăng ký học việc tại Thẩm mỹ viện Hà Anh do bà Tô Thị Ngân Hà làm chủ tại địa chỉ cũ 84 Đặng Nguyên Cẩn (hiện đã chuyển tới 110/20/1 đường Bà Hom, phường 13, quận 6, TP.HCM).
Do thiếu kinh nghiệm và không có kiến thức chuyên môn cũng như bằng cấp nên bà Hà đã tiêm sai kỹ thuật, bỏ qua các bước kiểm tra cơ bản và gây hậu quả là chị bị đột quỵ ngay khi vừa được tiêm filler, mắt trái hoàn toàn mất thị lực.
Cảnh nạn nhân bị tiêm vào mũi ( ảnh cắt từ clip nạn nhân cung cấp) Thấy tôi choáng váng và ói liên tục nên bà Hà đã nhờ em gái pha nước đường gừng cho tôi uống, nhưng chưa kịp uống thì tôi lại tiếp tục ói và gục xuống, lúc này bà Hà có gọi chồng về và kêu taxi đưa tôi đi cấp cứu.
Đêm hôm đó gia đình bà Hà về nhà, chỉ có bạn của tôi ở lại với tôi tại phòng cấp cứu của Bệnh viện 115 Bác sĩ Nguyễn Huy Thắng Trưởng khoa Bệnh lý Mạch máu não nói với ba mẹ tôi thì trong quá trình tiêm chất làm đầy cho tôi, bà Hà đã để mũi kim đâm trúng mạch máu làm thuyên tắc động mạch dẫn tới tôi bị xuất huyết não, yếu nửa người bên phải và mất thị lực mắt trái Một giáo sư chuyên khoa phẫu thuật tạo hình ở Hà Nội cho biết: Tiêm filler có thể gây phản ứng tại chỗ (phù, đau, ngứa, đỏ và sần da, có thể nhiễm trùng, áp xe); có thể gây tắc mạch làm mù mắt, tắc mạch não.
Các bác sĩ cũng kết luận là mắt chị Loan hoàn toàn không thể hồi phục được và rất khó tìm hướng điều trị tại Việt Nam vì toàn bộ dây thần kinh hốc mắt đã bị phá hủy kèm theo viêm màng bồ đào, viêm mống thể mi và bắt đầu bong giác mạc (mắt bị tổn thương nghiêm trọng không còn nguyên vẹn) do filler đi vào mạch máu.
Hiện tại chị Loan mới trải qua hai cuộc phẫu thuật tại Bệnh viện Mắt Cao Thắng TP.HCM để duy trì hình dạng nhãn cầu nhưng mắt vẫn trong tình trạng mất thị lực vĩnh viễn không có cách nào khôi phục được và phải điều trị để hạn chế teo mắt và giảm bớt các cơn đau.
Đơn khiếu nại của chị Loan gửi đến cơ quan chức năng Trong đơn khởi kiện gửi TAND quận 6, chị Loan yêu cầu bà Hà bồi thường 360 triệu đồng, gồm tiền bồi thường về tổn hại sức khỏe và chi phí điều trị khắc phục tình trạng teo nhãn cầu.
Điều 37 Nghị định 109/2016/NĐ-CP quy định về cấp chứng chỉ hành nghề đối với người hành nghề và hoạt động đối với cơ sở khám, chữa bệnh, có nêu: Cơ sở dịch vụ thẩm mỹ chỉ được thực hiện các hoạt động xăm, phun, thêu trên da, không sử dụng thuốc gây tê dạng tiêm Trường hợp này chủ cơ sở thẩm mỹ đã thực hiện tiêm vào cơ thể nạn nhân khi thực hiện quy trình thẩm mỹ; dạy nghề thẩm mỹ vi phạm quy định về khám, chữa bệnh do tiêm filler tại thẩm mỹ viện không có giấy phép hoạt động dưới hình thức BV và phòng khám chuyên khoa.
Theo luật sư Hưng, theo quy định tại Điều 242 bộ luật Hình sự năm 1999 (sửa đổi, bổ sung 2009) về tội vi phạm quy định về khám, chữa bệnh, sản xuất, pha chế thuốc, cấp phát thuốc, bán thuốc hoặc dịch vụ y tế khác thì: Người nào vi phạm quy định về khám, chữa bệnh, sản xuất, pha chế, cấp phát thuốc, bán thuốc hoặc dịch vụ y tế khác, nếu không thuộc trường hợp quy định tại điều 201 của bộ luật này, gây thiệt hại cho tính mạng hoặc gây thiệt hại nghiêm trọng cho sức khỏe của người khác hoặc đã bị xử lý kỷ luật, xử phạt hành chính về hành vi này hoặc đã bị kết án về tội này, chưa được xóa án tích mà còn vi phạm, thì bị phạt tù từ 1 5 năm.
Sau đó, chị L. đến khám tại Bệnh viện Chợ Rẫy, bác sĩ kết luận mắt trái của chị L. không thể hồi phục được vì toàn bộ dây thần kinh hốc mắt đã bị phá hủy, kèm theo viêm màng bồ đào, viêm mống thể mi và bắt đầu bong giác mạc.
Tuy nhiên quá trình tiêm đòi hỏi phải đúng chỉ định, bác sĩ phải có tay nghề cao, biết cấu trúc mạch máu để tiêm đúng vị trí, đúng cách, sử dụng loại mũi kim phù hợp, sản phẩm làm đầy phải đảm bảo chất lượng, thực hiện ở cơ sở y tế được cấp phép.
Và từ đó, dân gian mới có câu ca rằng: "Thế gian một vợ một chồng Chẳng như vua bếp hai ông một bà" Cả 3 tích truyện tuy có nhân vật hoàn cảnh khác nhau, nhưng đều có một điểm chung là những nhân vật đều sống có nghĩa có tình.
Qua câu chuyện trên, chúng ta thấy muốn được quả báo giàu sang sung sướng là do nhân bố thí đời trước, được quả báo thông minh là do nhân khuyên người khác làm lành tránh ác, quả báo tướng mạo đoan trang đẹp đẽ là do nhân đời trước giúp đỡ kẻ tật nguyền.
	`

	rawText = normalize(rawText)
	words := wordRe.FindAllString(rawText, -1)

	f, err := os.Create(OutputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	count := 0

	for win := MinWin; win <= MaxWin; win++ {
		for i := 0; i+win <= len(words); i++ {
			if count >= MaxPaths {
				fmt.Println("Generated:", count)
				return
			}

			phrase := strings.Join(words[i:i+win], " ")

			if utf8.RuneCountInString(phrase) > 150 {
				continue
			}

			hash := shortHash(phrase)
			tokens := strings.Split(phrase, "_")
			first := tokens[0]

			bucket := alphaBucket(first)

			path := fmt.Sprintf(
				"corpus/%s/%s/%s_%s.txt",
				bucket,
				first,
				phrase,
				hash,
			)

			w.WriteString(path + "\n")
			count++
		}
	}

	fmt.Println("Generated:", count)
}
