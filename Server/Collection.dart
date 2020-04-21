import 'dart:io';
import 'dart:convert';
import 'DataObject.dart';

class Collection{
  File file;
  String name;
  List<DataObject> dataList;


  Collection(String name){
    this.name = name;
    this.file = File("Databases/$name.json");
    dataList = new List<DataObject>();
    readFile();
    
  }

  void readFile() async {
    String content;
    if(await file.exists()){
      content = await file.readAsString();
      List<dynamic> newList = json.decode(content);
      for(var i in newList){
        DataObject obje = new DataObject.withMap(i["id"], i["data"]);
        dataList.add(obje);
      }
    }
  }
  void updateFile() async{
    if(await file.exists()){
      file.writeAsString(json.encode(dataList));
    }
  }
}