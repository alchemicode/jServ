import 'dart:io';
import 'dart:convert';
import 'DataObject.dart';

class Collection{
  File file;

  List<DataObject> dataList;

  Map map;

  Collection(String filename){
    file = File("Database/$filename");

  }
}