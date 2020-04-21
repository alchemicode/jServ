<p align="center">
    <img src="Media/icon.png" width="256px" height="256px">
</p>

<h1 align="center">
    jServ
</h1>

<p align="center">
    <img src="https://img.shields.io/github/license/Alchemi/jServ?style=flat-square">
    <img src="https://img.shields.io/github/manifest-json/v/Alchemi/jServ?style=flat-square">
    <img src="https://img.shields.io/badge/Build-In%20Development-red?style=flat-square">
    <img src="https://img.shields.io/badge/Platforms-Windows-brightgreen?style=flat-square">
    <a href="https://www.getpostman.com/collections/289f0bfba5cf1a9572c7">
        <img src="https://img.shields.io/badge/Postman-API_Ready-orange?style=flat-square&logo=postman">
    </a>
</p>


<h3 align="center">
    A flexible database server
</h3>
<br><br>
<table border="0">
    <tr>
        <th align="center">
            Build your backend fast, no strings attached
        </th> 
        <th align="center">
            Adapt the server's workflow to your project
        </th>
        <th align="center">
            Multiplatform, Simple, and easy to use
        </th>
        <th align="center">
            Use the built in API to run your apps
        </th>  
    </tr>
    <tr>
        <td align="center">
            jServ is an open source project designed to help backend developers get a server, database, and API up and running as soon as possible.<br>
        </td>
        <td align="center">
            jServ has a flexible data structure that allows you to customize the database and it's functionality, with or without modifying the code.<br>
        </td>
        <td align="center">
            jServ's driving force is it's use of JSON in it's data structure, allowing for a practical and effortless experience.<br>
        </td>
        <td align="center">
            jServ's API has been specifically designed to remain versatile and adaptable to your needs.<br>
        </td>
    </tr>
</table>
<br><br>

##Program Reference


jServ is extremely flexible. There are very few definite terms provided, as most of it depends on the implementation by the individual developer.


###Data Structure


The data structure relies on two classes, `DataObject` and `Collection`. 
 
`DataObject` is the class that all instances in the database come from. When serialized as a JSON object, it appears as the following,
```json
{
    "id": some-int,
    "data": {"some-key": some-value, ...}
}
```
The reason the object has only two attributes is that the developer defines what data each object will have within the `data` field. The `id` field is the only definite attribute to any object, as it is required for the API to be functional. It is dependent on the developer to ensure that the data field is consistent across all objects(if this is what is desired).
 
 
'Collection' is simply a container within the program for a database and its name. When written as a JSON object, it appears as the following,

```json
{
    "name": some-string,
    "dataList": [
        {
        "id": some-int,
        "data": {"some-key": some-value, ...}
        }
    ]
}
```
The `Collection` class exists to keep track of each database within the server. Within the program, the name corresponds to a filename in the `Databases` folder, which is what comprises the `dataList` in the class.


 ###API Reference

