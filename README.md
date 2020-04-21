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

<h2>Program Reference</h2>


jServ is extremely flexible. There are very few definite terms provided, as most of it depends on the implementation by the individual developer.


<h3>Data Structure</h3>


The data structure relies on two classes, `DataObject` and `Collection`. 
 

`DataObject` is the class that all instances in the database come from. When serialized as a JSON object, it appears as the following (with example values),
```json
{
    "id": 0,
    "data": {"some-key": "some-value"}
}
```

The reason the object has only two attributes is that the developer defines what data each object will have within the `data` field. The `id` field is the only definite attribute to any object, as it is required for the API to be functional. It is dependent on the developer to ensure that the data field is consistent across all objects(if this is what is desired).
 
<br>

`Collection` is simply a container within the program for a database and its name. When written as a JSON object, it appears as the following (with example values),

```json
{
    "name": "some-string",
    "dataList": [
        {
        "id": 0,
        "data": {"some-key": "some-value"}
        }
    ]
}
```

The `Collection` class exists to keep track of each database within the server. Within the program, the name corresponds to a filename in the `Databases` folder, which is what comprises the `dataList` in the class.


<h3>API Reference</h3>


jServ's API is built around a system of specific requests and query parameters.


<h4>GET Requests</h4>
 
<dl>
    <dt><code>__/query</code></dt>
    <dd>
    Queries a database for a specific object by id. Returns the whole object.
    <br>
    Query Parameters:
        <ul>
            <li>q - The name of the database you're querying
            <li>id - The id of the object you're querying
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/query/attribute</code></dt>
    <dd>
    Queries a database for a specific attribute of an object by id and name. Returns the attribute value.
    <br>
    Query Parameters:
        <ul>
            <li>q - The name of the database you're querying
            <li>id - The id of the object you're querying
            <li>a - the name of the attribute you're querying
        </ul>
    </dd>
</dl>
 
<h4>POST Requests</h4>

<dl>
    <dt><code>__/add</code></dt>
    <dd>
    Adds a new object to a database by id.
    <br>
    Query Parameters:
        <ul>
            <li>q - The name of the database you're adding to
            <li>id - The id of the object you're adding
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/add/data</code></dt>
    <dd>
    Adds an attribute to an object in a database by id.
    <br>
    Query Parameters:
        <ul>
            <li>q - The name of the database you're object is in
            <li>id - The id of the object you're adding to
            <li>a - The name of the attribute you're adding
            <li>v - The value of the attribute you're adding
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/mod/data</code></dt>
    <dd>
    Modifies an attribute of an object in a database by id.
    <br>
    Query Parameters:
        <ul>
            <li>q - The name of the database the object is in
            <li>id - The id of the object you're modifying to
            <li>a - The name of the attribute you're modifying
            <li>v - The new value of the attribute you're modifying
        </ul>
    </dd>
</dl>

