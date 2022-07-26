<p align="center">
    <img src="https://codealchemi.com/images/jserv.png" width="256px" height="256px">
</p>

<h1 align="center">
    jServ
</h1>

<h3 align="center">
    A project by codealchemi
</h3>

<p align="center">
    <img src="https://img.shields.io/badge/license-Apache%202.0-green?style=flat-square">
    <img src="https://img.shields.io/github/manifest-json/v/kketg/jServ?style=flat-square">
    <img src="https://img.shields.io/badge/Build-Functional-orange?style=flat-square">
    <img src="https://img.shields.io/badge/Platforms-Linux-brightgreen?style=flat-square">
</p>

<hr>

<h2 align="center">
    A flexible database server
</h2>

<br>
<table border="0">
    <tr>
        <th align="center">
            Build your backend fast, no strings attached
        </th> 
        <th align="center">
            Adapt the server's API and workflow to your project
        </th>
        <th align="center">
            Multiplatform, Simple, well-documented and easy to use
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
    </tr>
</table>
<br><br>

<h2>Getting Started</h2>
To set up jServ, download the latest release, and unzip it into a folder. You'll find the executable, a <code>config.json</code> file, a <code>version.jserv</code> file, a <code>keys.jserv</code> file, and a <code>data.jserv</code> file. 
There will also be a directory called <code>Databases</code>, with a <code>db.json</code> given to get started. To add a collection to the program, simply add a json file of any name, and add <code>[]</code> to the first line, and the program will read it.

Before you execute the program for the first time, you should check in your config and data files.

The <code>config.json</code> file should look something like this:
```json
{ 
    "ip": "localhost", 
    "port": 4040, 
    "Requests": { 
        "GET": true, 
        "POST": true, 
        "PUT": false, 
        "HEAD": true, 
        "DELETE": true, 
        "PATCH": false, 
        "OPTIONS": false 
    }, 
    "Permissions": { 
        "QObject": "user", 
        "QAttribute": "user", 
        "QAllAttributes": "user", 
        "QByAttribute": "user", 
        "QnewId": "admin", 
        "AEmpty": "user", 
        "AObject": "user", 
        "AAttribute": "user", 
        "MObject": "user", 
        "MAttribute": "user", 
        "DObject": "user", 
        "DAttribute": "user" 
    } 
}
```

Change the IP and port to whatever you desire. The requests list determines which requests the program will accept. For now, you can leave this alone.

The `keys.jserv` file will be empty. Do not manually add anything to the file.

The `data.jserv` file should look like this:
```
New App
new
```

Replace `New App` with the name of the application you are using jServ to develop.
DO NOT edit the `new` line, as that is where the program will generate the API key.

When you run the program, an Admin API key will generate in the `data.jserv` file, and a User API key will generate in the `keys.jserv` file.
The program will reject any requests that do not have these keys in the `"x-api-key"` header.


<h2>Program Reference</h2>


jServ is extremely flexible. There are very few definite terms provided, as most of it depends on the implementation by the individual developer.


<h3>Data Structure</h3>


The data structure relies on three classes, `DataObject`, `AttributeContainer`, and `Collection`. 
 

`DataObject` is the class that all instances in the database come from. When serialized as a JSON object, it appears as the following (with example values),
```json
{
    "id": 0,
    "data": {"some-key": "some-value"}
}
```

The reason the object has only two fields is that the developer defines what attributes each object will have within the `data` field. The `id` field is the only definite field to any object, as it is required for the API to be functional. It is dependent on you to implement field enforcement in your applications, and to ensure that the data fields are consistent across all objects(if this is what is desired). 
 
<br>

`AttributeContainer` is a class that serves the sole purpose of being a proxy between JSON objects passed in the API requests. When serialized as a JSON object, it appears as the following (with example values),
```json
{
    "some-key": "some-value"
}
```

Some of the requests require a single value to be passed in to the request body in the form of an `AttributeContainer` object, as this is the only way to maintain flexible typing within the database. The `AttributeContainer` class acts as a model within the program to translate that data seamlessly to the `Collection` and `DataObject` classes.
 
<br>

`Collection` is simply a container within the program for a database and its name. When written as a JSON object, it appears as the following (with example values),

```json
{
    "name": "some-string",
    "data-list": [
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
    Queries a collection for a specific object by id. Returns the whole object in JSON.
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're querying</li>
            <li>id - The id of the object you're querying</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/query/attribute</code></dt>
    <dd>
    Queries a collection for a specific attribute of an object by id and name. Returns the attribute value in an <code>AttributeContainer</code> object.
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're querying</li>
            <li>id - The id of the object you're querying</li>
            <li>a - The name of the attribute you're querying</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/query/allAttributes</code></dt>
    <dd>
    Queries a collection for all attributes of a specific key in every object. If an object does not have an attribute of the passed key, the object is skipped. The query returns a list of all the attributes keyed by object id.
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're querying</li>
            <li>a - the name of the attributes you're querying</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/query/byAttribute</code></dt>
    <dd>
    Queries a collection for objects that share the same value of a specific attribute. If an object does not have an attribute of the passed key, the object is skipped. The query returns a list of all the objects with the attribute and value. (<em>Requires an <code>AttributeContainer</code> JSON object to be passed in the body</em>)
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're querying</li>
            <li>a - The name of the attributes you're querying</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/query/newId</code></dt>
    <dd>
    Returns an unused id in a collection
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're querying</li>
        </ul>
    </dd>
</dl>
 
<h4>POST Requests</h4>

<dl>
    <dt><code>__/add</code></dt>
    <dd>
    Adds a new empty object to a collection by id.
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're adding to</li>
            <li>id - The id of the object you're adding</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/add/object</code></dt>
    <dd>
        Adds a new JSON object to a collection (<em>Requires an <code>DataObject</code> JSON object to be passed in the body</em>).
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're adding to</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/add/attribute</code></dt>
    <dd>
    Adds an attribute to an object in a collection by id (<em>Requires an <code>AttributeContainer</code> JSON object to be passed in the body</em>).
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're object is in</li>
            <li>id - The id of the object you're adding to</li>
            <li>a - The name of the attribute you're adding</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/mod/object</code></dt>
    <dd>
    Modifies the id of an object in a collection by id.
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection the object is in</li>
            <li>id - The id of the object you're modifying</li>
            <li>v - The new id of the object you're modifying</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/mod/attribute</code></dt>
    <dd>
    Modifies an attribute of an object in a collection by id (<em>Requires an <code>AttributeContainer</code> JSON object to be passed in the body</em>).
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection the object is in</li>
            <li>id - The id of the object you're modifying</li>
            <li>a - The name of the attribute you're modifying</li>
        </ul>
    </dd>
</dl>

<h4>DELETE Requests</h4>

<dl>
    <dt><code>__/delete/object</code></dt>
    <dd>
    Deletes an object from a collection by id.
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're deleting from</li>
            <li>id - The id of the object you're deleting</li>
        </ul>
    </dd>
</dl>
<dl>
    <dt><code>__/delete/attribute</code></dt>
    <dd>
        Deletes an attribute from an object by id.
    <br>
    Query Parameters:
        <ul>
            <li>db - The name of the collection you're deleting from</li>
            <li>id - The id of the object you're deleting</li>
            <li>a - The name of the attribute you're deleting</li>
        </ul>
    </dd>
</dl>
<br>
<h2 align="center">License and Copyright Notice</h2>
<p align="center">
    Copyright (c) 2022, Kristofer Ter-Gabrielyan. All Rights Reserved. 
    Permission to modify and redistribute is granted under the terms of the Apache 2.0.
</p>
