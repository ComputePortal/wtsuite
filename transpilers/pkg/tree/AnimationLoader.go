package tree

import (
	"strconv"
	"strings"
)

// this is function is guaranteed to be called after document is fully loaded
func WriteAnimationLoader(scenes []int) string {
	var b strings.Builder

	// get cookie utility function
	b.WriteString("function gc(n){")
	b.WriteString("n=n+\"=\";")
	b.WriteString("let d=decodeURIComponent(document.cookie);")
	b.WriteString("let l=d.split(';');")
	b.WriteString("for(let i=0;i<l.length;i++){")
	b.WriteString("let c=l[i];")
	b.WriteString("while(c.charAt(0)==' '){c=c.substring(1);}")
	b.WriteString("if(c.indexOf(n)==0){return c.substring(n.length,c.length);}")
	b.WriteString("}")
	b.WriteString("return \"\";}")

	// set cookie
	b.WriteString("function sc(n,v){")
	b.WriteString("let t=new Date();t.setTime(t.getTime()+1000*3600);")
	b.WriteString("document.cookie=n+'='+v+';expires='+t.toGMTString()+';path=/;';")
	b.WriteString("}")

	// list of scenes as variable
	b.WriteString("let scenes=[")
	for i, val := range scenes {

		b.WriteString(strconv.Itoa(val))

		if i < len(scenes)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("];")

	// keydown isnt able to detect PrintScreen in Windows, so must use keyup
	b.WriteString("window.addEventListener('keyup', (e)=>{if(e.key=='PrintScreen'){")
	b.WriteString("let c=gc('scene');")
	b.WriteString("if(c==''){")
	// set the cookie to prev and refresh the page
	b.WriteString("sc('scene', scenes[0]-1);")
	b.WriteString("window.location.href=window.location.href;")
	b.WriteString("}else{")
	b.WriteString("document.cookie='scene=;path=/';")
	// stop the animation by refreshing the page
	b.WriteString("window.location.href=window.location.href;")
	b.WriteString("}}});")

	// for check if cookie exists
	b.WriteString("let fi=-1;let c;let cprev=gc('scene');")
	b.WriteString("if(cprev!==undefined&&cprev!=''){")
	b.WriteString("c=parseInt(cprev)+1;")
	b.WriteString("fi=scenes.findIndex((s)=>s==c);")
	b.WriteString("}")

	b.WriteString("if(fi!=-1){")
	// the animation can be continued
	// update the cookie immediately
	b.WriteString("sc('scene',c);")
	// TODO: custom scene file prefix
	b.WriteString("import('/scene'+c.toString()+'.js');")
	b.WriteString("}")

	return b.String()
}
