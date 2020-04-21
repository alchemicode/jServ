import 'dart:convert';

class AttributeContainer{
  Map<String, dynamic> data = new Map<String, dynamic>();
  String key;
  dynamic value;
  AttributeContainer(String key, dynamic value){
    this.key = key;
    this.value = value;
    data[key] = value;
  }

  String toJSON(){
    return json.encode(data);
  }

  String toString(){
    return "{ \"" + key + "\" : " + value.toString() + " }";
  }

}
