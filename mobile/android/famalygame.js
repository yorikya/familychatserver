var _clientid ="yuri@gmail.com"
var _serverip = "10.0.0.2:8080"
var _client_data =""


function SendLoginRequest(host, clientid) {
 //Send request to remote server.
    var url = "http://"+host;
    var path = "/auth";
    var params = "roomid=1|id="+clientid;
    app.HttpRequest( "get", url, path, params, handleReply );   
}

function handleReply( error, reply )
{
    if( error ) alert( error );
    else
    {
        alert( reply );
        _client_data = JSON.parse(reply)
    }
}

//Handle download completion.
function dload_OnDownload( file )
{
	app.ShowPopup( "Downloaded: " + file );
		//Create image 1/5 of screen width and correct aspect ratio.
// 	img = app.CreateImage( "/sdcard/temp/superman.jpeg", 0.15, 0.1 );
// 	lay.AddChild( img );
}

//Handle download errors.
function dload_OnError( error )
{
	app.ShowPopup( "Download failed: " + error );
}


//Called when application is started.
function OnStart()
{   
    
    
    //  var url = "http://"+_serverip+"/resources/rooms/1/superman.jpeg,http://"+_serverip+"/resources/rooms/1/wonderwoman.jpeg";
    //  app.Debug("the download url"+url)
    //  var fldr = "/sdcard/DroidScript/resources/rooms/1";
    
    // //Make sure target folder exists.
    // app.MakeFolder( fldr );
    
    // //Download file from web.
    // //(You can leave out the file parameter to use original file name)
    // dload = app.CreateDownloader(/*"NoDialog"*/  );
    // dload.SetOnDownload( dload_OnDownload );
    // dload.SetOnError( dload_OnError );
    // dload.Download( url, fldr );
    
 

    //Create and run web server. 
	serv = app.CreateWebServer( 8080, "Upload,ListDir" ); 
	serv.SetFolder( "/sdcard/DroidScript" ); 
	serv.AddServlet( "/message", OnServlet ); 
	serv.Start(); 
    
    //Create the login page.
	lay = app.CreateLayout( "Linear", "FillXY" );	
	lay.SetBackground( "/Sys/Img/BlueBack.jpg" );
	lay.SetPadding( 0, 0.1, 0, 0 ); 

    username = app.CreateText( "User Name" , 0.3, 0.04);
    username.SetMargins(-0.2,0,0,0)
    username.SetTextSize( 20 );
    lay.AddChild( username );
    
	userID = app.CreateTextEdit( "", 0.7, 0.1 );
	userID.SetMargins(0,0,0,0.1)
    lay.AddChild( userID );
    
    passlable = app.CreateText( "Password" , 0.3, 0.04);
    passlable.SetMargins(-0.2,0,0,0)
    passlable.SetTextSize( 20 );
    lay.AddChild( passlable );
    
    password = app.CreateTextEdit( "", 0.7, 0.1);
    lay.AddChild( password );
    
    //Create button and add to main layout.
	loginBtn = app.CreateButton( "Login", 0.4, 0.1, "gray" );
	loginBtn.SetMargins(0,0.1,0,0)
	loginBtn.SetOnTouch( loginBtnOnTouch );
	lay.AddChild( loginBtn );
    

	//Create a layout we can slide over the main layout.
	//(This hidden layout is actually on top of the main
	//layout, it just appears to slide from the left)
	laySlide = app.CreateLayout( "Linear", "FillXY" );
	laySlide.SetPadding( 0, 0.1, 0, 0 ); 
	laySlide.SetBackground( "/Sys/Img/GreenBack.jpg" );
	laySlide.SetVisibility( "Hide" );
	
	//Create button and add to sliding layout.
	btnBack = app.CreateButton( "Back", 0.3, 0.06, "gray" );
	btnBack.SetOnTouch( btnBack_OnTouch );
	laySlide.AddChild( btnBack );
	
// 	//Create a layout with objects vertically centered.
// 	lay = app.CreateLayout( "linear", "VCenter,FillXY" );	

	txt = app.CreateList("Conversations", 0.8, 0.4  );
	laySlide.AddChild( txt );
	//Create a button 1/3 of screen width and 1/10 screen height.
	btn = app.CreateButton( "Press Me", 0.3, 0.1 );
	btn.SetMargins( 0, 0.05, 0, 0 );
	
	//Set function to call when button pressed.
	btn.SetOnTouch( btn_OnTouch );
	laySlide.AddChild( btn );
	
	edt = app.CreateTextEdit( "Hello", 0.7, 0.1, "Multiline" );
    edt.SetTextColor( "#ff6666ff" );
    edt.SetBackColor( "#ffffffff" );
    laySlide.AddChild( edt );
    
    //Create a text label and add it to layout.
	txtMsg = app.CreateText( "", 0.8, 0.3, "AutoScale,MultiLine" );
	txtMsg.SetTextSize( 22 );
	laySlide.AddChild( txtMsg );

	//Add layout to app.	
	app.AddLayout( lay );
	app.AddLayout( laySlide );
	
// 	setInterval( ShowServerResource, 1000 );
}

function ShowServerResource() {
     app.Debug("update resource path to:" + _client_data.resources);
}

//Handle servlet requests. 
function OnServlet( request, info ) 
{ 
	serv.SetResponse( "Got it!" ); 
	app.Debug("response text", request.msg );
    txt.AddItem(request.id+ ": "+ request.msg);
} 
 
//Send an http get request. 
function SendRequest( url ) 
{ 
    var httpRequest = new XMLHttpRequest(); 
    httpRequest.onreadystatechange = function() { HandleReply(httpRequest); };   
    httpRequest.open("GET", url, true); 
    httpRequest.send(null); 
 
    app.ShowProgress( "Loading..." ); 
} 
 

//Handle the server's reply (a json object). 
function HandleReply( httpRequest ) 
{ 
    if( httpRequest.readyState==4 ) 
    { 
        //An error occurred  
        if( httpRequest.status != 200 ) 
        { 
            app.Alert( "Error: " + httpRequest.status + httpRequest.responseText);
        }  
    } 
  app.HideProgress(); 
} 


//Called when user touches our slide button.
function loginBtnOnTouch()
{
	laySlide.Animate( "SlideFromLeft" );
}

//Called when user touches our back button.
function btnBack_OnTouch()
{
	laySlide.Animate( "SlideToLeft" );	
}

//Called when user touches our button.
function btn_OnTouch()
{
    SendRequest("http://"+_serverip+"/broadcast?id="+_clientid+"&msg="+ edt.GetText());
    edt.SetText("");
}