var _serverhost = "http://10.0.0.2:8080"
var _client_data;
var _app_fldr = "/sdcard/DroidScript"
var _app_login_lay;
var _app_chat_lay;

function createLoginLayout() {
      //Create the login page.
	lay = app.CreateLayout( "Linear", "FillXY" );	
	lay.SetBackground( "/Sys/Img/BlueBack.jpg" );
	lay.SetPadding( 0, 0.1, 0, 0 ); 

    username = app.CreateText( "User Name" , 0.3, 0.04);
    username.SetMargins(-0.2,0,0,0)
    username.SetTextSize( 20 );
    lay.AddChild( username );
    
	userID = app.CreateTextEdit( "yuri", 0.7, 0.1 );
	userID.SetMargins(0,0,0,0.1)
    lay.AddChild( userID );
    
    passlable = app.CreateText( "Password" , 0.3, 0.04);
    passlable.SetMargins(-0.2,0,0,0)
    passlable.SetTextSize( 20 );
    lay.AddChild( passlable );
    
    password = app.CreateTextEdit( "a", 0.7, 0.1);
    lay.AddChild( password );
    
    //Create button and add to main layout.
	loginBtn = app.CreateButton( "Login", 0.4, 0.1, "gray" );
	loginBtn.SetMargins(0,0.1,0,0)
	loginBtn.SetOnTouch( function(){sendLoginRequest(userID.GetText(), password.GetText());});
	lay.AddChild( loginBtn );
	return lay
}

function createChatLayout() {
	lay = app.CreateLayout( "Linear", "FillXY" );
	lay.SetPadding( 0, 0.1, 0, 0 ); 
	lay.SetBackground( "/Sys/Img/GreenBack.jpg" );
	lay.SetVisibility( "Hide" );
	
	//Create button and add to sliding layout.
	btnBack = app.CreateButton( "Back", 0.3, 0.06, "gray" );
	btnBack.SetOnTouch( btnBack_OnTouch );
	lay.AddChild( btnBack );
	
	txt = app.CreateList("Conversations", 0.8, 0.4  );
	lay.AddChild( txt );
	//Create a button 1/3 of screen width and 1/10 screen height.
	btn = app.CreateButton( "Press Me", 0.3, 0.1 );
	btn.SetMargins( 0, 0.05, 0, 0 );
	
	//Set function to call when button pressed.
	btn.SetOnTouch( btn_OnTouch );
	lay.AddChild( btn );
	
	edt = app.CreateTextEdit( "Hello", 0.7, 0.1, "Multiline" );
    edt.SetTextColor( "#ff6666ff" );
    edt.SetBackColor( "#ffffffff" );
    edt.SetOnChange(function() { 
        app.Debug( edt.GetText() );
    });
    lay.AddChild( edt );
    
    //Create a text label and add it to layout.
	txtMsg = app.CreateText( "", 0.8, 0.3, "AutoScale,MultiLine" );
	txtMsg.SetTextSize( 22 );
	lay.AddChild( txtMsg );
	return lay
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
	
// 	setInterval( ShowServerResource, 1000 );
}

function sendLoginRequest(user, pass) {
    //Send request to remote server.
    let path = "/auth";
    let params = "user="+user+"|pass="+pass;
    app.HttpRequest( "get", _serverhost, path, params, function ( error, reply )
    {
        if( error ) alert( error );
        else {
            r = JSON.parse(reply);
            alert(reply);
            if (r.success == true) {
                _client_data = r;
                let urls = [];
                _client_data.filenames.forEach(function (name, _) {
                    urls.push(_serverhost+_client_data.resources+name)
                });

                let fldr = _app_fldr + _client_data.resources;
                //Make sure target folder exists.
                app.MakeFolder( fldr );
                
                dload = app.CreateDownloader(/*"NoDialog"*/  );
                dload.SetOnDownload( function(file) { _app_chat_lay.Animate( "SlideFromLeft" );});
                dload.SetOnError( function(error) { alert(error);});
                dload.Download( urls, fldr );
                            
                } else {
                    alert( r.error)
                }
                
            }
        });   
}

//Handle servlet requests. 
function onServlet( request, info ) { 
	serv.SetResponse( "Got it!" ); 
// 	txt = app.CreateList("Conversations", 0.8, 0.4  );
    txt.AddItem(request.id+ ": "+ request.msg);
} 

//Called when user touches our back button.
function btnBack_OnTouch() {
	_app_chat_lay.Animate( "SlideToLeft" );	
}

function sendMessage(message) {
      //Send request to remote server.
    let path = "/broadcast";
    let params = "id="+_client_data.clientid+"|msg="+message;
    app.HttpRequest( "get", _serverhost, path, params, function ( error, reply ) {
        if( error ) alert( error );        
    });   
}

//Called when user touches our button.
function btn_OnTouch() {
    sendMessage(edt.GetText());
    edt.SetText("");
}