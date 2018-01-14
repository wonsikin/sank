package downloader

import "testing"

func TestDownloadFile(t *testing.T) {
	name, url := "xxx", "http://wx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0"
	fileName, err := DownloadFile("./images/", name, url)
	t.Logf("error is %s, file name is %s", err, fileName)
}
