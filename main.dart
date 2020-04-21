import 'dart:io';
import 'dart:convert';
import 'DataObject.dart';

Map<String, dynamic> requestTypes = {"GET":true,"POST":true,"PUT":true,"HEAD":false,"DELETE":false,"PATCH":false,"OPTIONS":false};
String ip = "localhost";
int port = 4040;


Future main() async {
  //readConfig();
  
  
  var server = await HttpServer.bind(
    InternetAddress.loopbackIPv4,
    port,
  );
  var obj = new DataObject.withMap(10,{"e": 1, "yee":"ooo"});
  print("Listening on localhost:${server.port}");

  await for (HttpRequest request in server) {
    String m = request.method;
    if(requestTypes[m]){
      handleRequest(request);
    }else{
      request.response.write("Error: Invalid request");
    }
    
  }

  
}
void handleRequest(HttpRequest r){
    String m = r.method;
    try{
      switch(m){
        case("GET"):{
          handleGet(r);
        }
        break;
        case("POST"):{
          handlePost(r);
        }
        break;
        case("PUT"):{
          handlePost(r);
        }
        break;
        case("HEAD"):{
          handlePost(r);
        }
        break;
        case("DELETE"):{
          handlePost(r);
        }
        break;
        case("PATCH"):{
          handlePost(r);
        }
        break;
        case("OPTIONS"):{
          handlePost(r);
        }
        break;
      }
    }catch(e){
      print("Exception in handleRequest: $e");
    }
}
void handleGet(HttpRequest r){
  String query = r.uri.queryParameters["id"];
  print(r.uri.queryParameters["id"]);
  var response = r.response;
  response.write(query);
  response.close();

}
void handlePost(HttpRequest r){
  
}
void handlePut(HttpRequest r){
  
}
void handleHead(HttpRequest r){
  
}
void handleDelete(HttpRequest r){
  String m = r.method;
  r.response.write("this is a $m");
}
void handleOptions(HttpRequest r){
  String m = r.method;
  r.response.write("this is a $m");
}

Future readConfig() async {
  File file = new File("config.json");
  String content;

  if(await file.exists()){
    content = await file.readAsString();
  }
  var map = json.decode(content);
  ip = map["ip"];
  port = map["port"];
  requestTypes = map["Requests"];
  map.forEach((k,v) => print("$k : $v"));
  requestTypes.forEach((k,v) => print("$k : $v"));
}





