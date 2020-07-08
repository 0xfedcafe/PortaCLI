package main

import "fmt"

var (
	startMessage = `
                                 __                 ______   __        ______ 
                                /  |               /      \ /  |      /      |
  ______    ______    ______   _$$ |_     ______  /$$$$$$  |$$ |      $$$$$$/ 
 /      \  /      \  /      \ / $$   |   /      \ $$ |  $$/ $$ |        $$ |  
/$$$$$$  |/$$$$$$  |/$$$$$$  |$$$$$$/    $$$$$$  |$$ |      $$ |        $$ |  
$$ |  $$ |$$ |  $$ |$$ |  $$/   $$ | __  /    $$ |$$ |   __ $$ |        $$ |  
$$ |__$$ |$$ \__$$ |$$ |        $$ |/  |/$$$$$$$ |$$ \__/  |$$ |_____  _$$ |_ 
$$    $$/ $$    $$/ $$ |        $$  $$/ $$    $$ |$$    $$/ $$       |/ $$   |
$$$$$$$/   $$$$$$/  $$/          $$$$/   $$$$$$$/  $$$$$$/  $$$$$$$$/ $$$$$$/ 
$$ |                                                                          
$$ |                                                                          
$$/
Thank you for using PortaCLI.
Authors:
	godande
	MaxUNof
Donate BTC: 36MQNEv8vkXgVuTa8HS1aJYzFTsuCmwNBK
`
	usage = func() {
		fmt.Println(`Usage:
	-f = Path to your photo/video                | Default: img.jpg in your current folder
	-o = Change output destination with filename | Default: img_portacli.jpg/img_portacli.mp4 in your current folder
    -vid = Create video portrait
    -collage = Collage mode (only with photo)
Example usages:
		portacli -f /home/user/myphoto.jpg -o /home/user/outphoto.jpg
		portacli -f /home/user/myphoto.jpg -vid`)
	}
)
