package util

import "testing"

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/7
 * 
*/

func TestMD5(t *testing.T) {
	output := MD5("Windear")
	if output != "fa8a0d32f3772356dfcf99a1a354dfe3" {
		t.Errorf("Md5 Wrong Output With:%v", output)
	}
}