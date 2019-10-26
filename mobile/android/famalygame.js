var _serverhost = "http://10.0.0.2:8080"
var _client_data;
var _app_fldr = "/sdcard/DroidScript";
var _app_login_lay;
var _app_chat_lay;

function createLoginLayout() {
      //Create the login page.
	let lay = app.CreateLayout( "Linear", "FillXY" );	
	lay.SetBackground( "/Sys/Img/BlueBack.jpg" );
	lay.SetPadding( 0, 0.1, 0, 0 ); 

    let login_userid_lbl = app.CreateText( "User Name" , 0.3, 0.04);
    login_userid_lbl.SetMargins(-0.2,0,0,0);
    login_userid_lbl.SetTextSize( 20 );
    lay.AddChild( login_userid_lbl );
    
	login_userid_edt = app.CreateTextEdit( "yuri", 0.7, 0.1 );
	login_userid_edt.SetMargins(0,0,0,0.1)
    lay.AddChild( login_userid_edt );
    
    let login_pass_lbl = app.CreateText( "Password" , 0.3, 0.04);
    login_pass_lbl.SetMargins(-0.2,0,0,0)
    login_pass_lbl.SetTextSize( 20 );
    lay.AddChild( login_pass_lbl );
    
    login_pass_edt = app.CreateTextEdit( "a", 0.7, 0.1);
    lay.AddChild( login_pass_edt );
    
    //Create button and add to main layout.
	let login_btn = app.CreateButton( "Login", 0.4, 0.1, "gray" );
	login_btn.SetMargins(0,0.1,0,0)
	login_btn.SetOnTouch( function(){sendLoginRequest(login_userid_edt.GetText(), login_pass_edt.GetText());});
	lay.AddChild( login_btn );
	
	return lay;
}

function createChatLayout() {
	let lay = app.CreateLayout( "Linear", "FillXY" );
	lay.SetPadding( 0, 0.1, 0, 0 ); 
	lay.SetBackground( "/Sys/Img/GreenBack.jpg" );
	lay.SetVisibility( "Hide" );
	
	//Create button and add to sliding layout.
	chat_back_btn = app.CreateButton( "Back", 0.3, 0.06, "gray" );
	chat_back_btn.SetOnTouch(function(){ _app_chat_lay.Animate( "SlideToLeft" ); });
	lay.AddChild( chat_back_btn );
	
	chat_converse_lst = app.CreateList("Conversations", 0.8, 0.4  );
	lay.AddChild( chat_converse_lst );
	
	chat_msg_send_lay = app.CreateLayout( "Linear", "Horizontal,FillXY" );
	
	chat_msg_edt = app.CreateTextEdit( "Hello", 0.75, 0.07, "Multiline" );
    chat_msg_edt.SetTextColor( "#ff6666ff" );
    chat_msg_edt.SetBackColor( "#ffffffff" );
    chat_msg_edt.SetOnChange(function() { 
        let t = chat_msg_edt.GetText();
        let last_char = t.substring(t.length-1, t.length);
        if (last_char == "\n") {
            sendMessage(t.replace(/\n/g, ''));
            chat_msg_edt.SetText("");  
            chat_converse_lst.ScrollToItemByIndex(chat_converse_lst.GetLength());
        }
        
    });
    chat_msg_send_lay.AddChild(chat_msg_edt);
    
	let chat_send_btn = app.CreateButton( "Send", 0.2, 0.07 );
	chat_send_btn.SetMargins( 0.015, 0, 0, 0 );
	chat_send_btn.SetOnTouch(function(){ 
	    sendMessage(chat_msg_edt.GetText());
        chat_msg_edt.SetText(""); 
    });    
	chat_msg_send_lay.AddChild( chat_send_btn );
	
	lay.AddChild(chat_msg_send_lay)
	
	return lay;
}

//Main: Called when application is started.
function OnStart() {   
    //Create and run web server. 
	serv = app.CreateWebServer( 8080, "Upload,ListDir" ); 
	serv.SetFolder( "/sdcard/DroidScript" ); 
	serv.AddServlet( "/message", onServlet ); 
	serv.Start(); 
    
    //Create the login page.
    _app_login_lay = createLoginLayout();
    _app_chat_lay = createChatLayout();

	//Add layout to app.	
	app.AddLayout( _app_login_lay );
	app.AddLayout( _app_chat_lay );
	
}

function sendLoginRequest(user, pass) {
    //Send request to remote server.
    let path = "/auth";
    let params = "user="+user+"|pass="+pass;
    app.HttpRequest( "get", _serverhost, path, params, function ( error, reply )
    {
        if( error ) alert( error );
        else {
            let r = JSON.parse(reply);
            alert(reply);
            if (r.success == true) {
                _client_data = r;
                let urls = [];
                _client_data.filenames.forEach(function (name, _) {
                    urls.push(_serverhost+_client_data.resources+name);
                });

                let fldr = _app_fldr + _client_data.resources;
                //Make sure target folder exists.
                app.MakeFolder( fldr );
                
                dload = app.CreateDownloader(/*"NoDialog"*/  );
                dload.SetOnError( function(error) { alert(error);});
                dload.Download( urls, fldr );
                
                _app_chat_lay.Animate( "SlideFromLeft" );
                
                } else {
                    alert( r.error);
                }
                
            }
        });   
}

//Handle servlet requests. 
function onServlet( request, info ) { 
	serv.SetResponse( "Got it!" ); 
    chat_converse_lst.AddItem(request.id, request.msg + "\n\n" + request.ts);
    chat_converse_lst.ScrollToItemByIndex(chat_converse_lst.GetLength());
} 

function sendMessage(message) {
      //Send request to remote server.
    let path = "/broadcast";
    let params = "id="+_client_data.clientid+"|msg="+message;
    app.HttpRequest( "get", _serverhost, path, params, function ( error, reply ) {
        if( error ) alert( error );        
    });   
}