import 'dart:convert';


class DataObject{ 
  int id;
  Map data;
  
  DataObject.emptyMap(id){
    this.id = id;
    data = new Map();
  }

  DataObject.withMap(this.id, this.data);

  DataObject.fromJsonMap(Map<String, dynamic> map){
    this.id = map["id"];
    this.data = map["data"];
  }
  DataObject.fromJsonString(String s){
    Map<String, dynamic> map = json.decode(s);
  }


  Map toJson(){
    Map map = new Map();
    map["id"] = this.id;
    map["data"] = this.data;
    return map;
  }

  String toString(){
    return "id: " + id.toString() + "\ndata: " + data.toString();

  }
  
}