<p align="center">
    <img src="https://alchemicode.com/images/logos/jserv.png" width="256px" height="256px">
</p>

<h1 align="center">
    jServ
</h1>

<h3 align="center">
    A project by alchemicode
</h3>

<p align="center">
    <img src="https://img.shields.io/badge/license-Apache%202.0-green?style=flat-square">
    <img src="https://img.shields.io/github/manifest-json/v/kketg/jServ?style=flat-square">
    <img src="https://img.shields.io/badge/Build-Functional-orange?style=flat-square">
    <img src="https://img.shields.io/badge/Platforms-Linux-brightgreen?style=flat-square">
</p>

<hr>

<h2 align="center">
    A flexible backend server
</h2>

<br>
<table border="0">
    <tr>
        <th align="center">
            Build your backend fast and simple
        </th> 
        <th align="center">
            Easy document-based database with no query language
        </th>
        <th align="center">
            Adapt the server to your project with broad markup support and custom endpoint functions
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
            jServ's driving force is it's use of multiple markup languages in it's data structure, and scriptable endpoints.<br>
        </td>
    </tr>
</table>
<br><br>

<h2>Getting Started</h2>
<p>
To set up jServ, download the latest release, and unzip it into a folder, and run the executable. You will have a <code>config.json</code> file, a <code>keys.jserv</code> file, and an <code>admin.jserv</code> file. 
There will also be a directory called <code>Collections</code>, with an <code>example.dat</code> given to get started. To add a collection to the program, simply add a <code>.dat</code> file of any name, and the program will read it.
</p>

<p>
Before you execute the program for the first time, you should check in your config and data files.
</p>

The <code>config.json</code> file should look something like this:
```json
{   
    "appname": "New app",
    "debug": true,
    "admin-path": "admin.jserv",
    "key-path": "keys.jserv",
    "ip": "localhost", 
    "port": 4040, 
    "write-interval": 10,
    "Requests": 
    {   
        "GET": true, 
        "POST": true, 
        "PUT": false, 
        "HEAD": true, 
        "DELETE": true, 
        "PATCH": false, 
        "OPTIONS": false 
    }, 
    "Permissions": 
    { 
        "Query":  "user",
        "Add":    "user",
        "Mod":    "user",
        "Delete": "user",
        "Purge":  "user"
    },
    "Aliases":
    {
        "127.0.0.1":"localhost"
    },
    "Services":
    {
        "/" : "index"
    }
}
```
<p>
Change the IP and port to whatever you need. Debug mode will show more detailed console logging. The requests list determines which requests the program will accept.

The permissions list determines which requests can be made with the user keys, whereas admin keys will have access to all of them.

The aliases list will change how certain addresses are displayed within the console.

When you run the program for the first time, an Admin API key will generate in the <code>admin.jserv</code> file, and a User API key will generate in the <code>keys.jserv</code> file.
The program will reject any requests that do not have these keys in the <code>"x-api-key"</code> header.
</p>

<h1>Program Reference</h1>

<h2>Database</h2>

<p>
jServ's database relies on the use of HTTP requests to send instructions and data back and forth between the instance and your program.
There are built-in request handlers that execute a variety of database operations (See below).

Each request will give a response in the form of a JSON object. It appears as follows (with example values)
```json
{
    "status": "ok",
    "message": "Successfully queried some-database for some-object",
    "data": {
        "some-data": "some-value"
    }
}
```
<p>
The <code>status</code> value will appear as either <code>"ok"</code> or <code>"error"</code>, and the <code>message</code> value will display a message either confirming the success, or explaining the error. The <code>data</code> value may not appear, and will only contain data if the request returns it.
</p>

<h3>Data Structure</h3>
(All of these structures can also be found in the <a href="https://github.com/alchemicode/jserv-core">jServ core library</a>)
<br><br>
<p>
The database follows a document-based structure which internally relies on three classes, <code>Document</code>, <code>Attribute</code>, and <code>Collection</code>. 
</p>

<p>
<code>Document</code> is the class that represents objects in the database. When serialized as a JSON object, it appears as follows (with example values)
</p>

```json
{
    "_id": "some-unique-identifier",
    "data": {"some-key": "some-value"}
}
```
<p>
There is no pre-defined schema like in relational databases like MySQL. Each document can have a unique set of <code>data</code> The `_id` field is the only guaranteed value attached to any document. It is dependent on the user to implement field enforcement in your applications, and to ensure that the data fields are consistent across all objects if that is required. 
</p>
<br>
<p>
<code>Attribute</code> is a class that serves the sole purpose of being a proxy between fields passed in the API requests. When serialized as a JSON object, it appears as follows (with example values)
</p>

```json
{
    "some-key": "some-value"
}
```
<p>
Some of the requests require a single value to be passed in to the request body in a form resembling an <code>Attribute</code> object. The <code>Attribute</code> class acts as a model within the program to translate that data seamlessly to the <code>Collection</code> and <code>Document</code> classes.
</p>
<br>
<p>
<code>Collection</code> is a grouping of documents within a file. When written as a JSON object, it appears as follows (with example values)
</p>

```json
{
    "name": "some-name",
    "list": [
        {
        "id": 0,
        "data": {"some-key": "some-value"}
        }
    ]
}
```
<p>
Internally, the name corresponds to a filename in the <code>Collections</code> folder, which contains the serialized data of <code>list</code>.
</p>

<br>

<h2>Database API Reference</h2>
(All of these structures can also be found in the <a href="https://github.com/alchemicode/jserv-core">jServ core library</a>)
<br><br>
jServ's database operations are called through HTTP requests. This eliminates the need for a query language, as the properties of every operation can be encapsulated into serialized objects passed into each request.

<h3>Operations</h3>
<dl>
<dt><code>__/j/db/query</code></dt>
<dd>
<p>
Queries a list of, or all, collections for documents based on a set of rules. Returns a list of documents.
<br>
Possible properties of a query, written in JSON format, are as follows:
</p>
        
```json
{
    //List of collections to be queried
    //If omitted, all collections will be queried
    "collections": [
        "some-collection", 
        "another-collection"
    ],
    //List of attributes a document must have
    "has": [
        "some-attribute",
        "another-attribute"
    ],
    //Attributes that must be equal to a given value
    "equals": {
        "some-attribute": "some-value"
    },
    //Attributes that must not be equal to a given value
    "not-equals": {
        "some-attribute": "some-value"
    },
    //Attributes that must be less than a given value
    "lt": {
        "some-attribute": "some-value"
    },
    //Attributes that must be less than or equal to a given value
    "lte": {
        "some-attribute": "some-value"
    },
    //Attributes that must be greater than or equal to a given value
    "gte": {
        "some-attribute": "some-value"
    },
    //Attributes that must be greater than a given value
    "gt": {
        "some-attribute": "some-value"
    },
    //Attributes that must be between two given values
    "between": {
        "some-attribute": ["lower-value", "upper-value"]
    }
}
```
<p>This request returns the list of documents queried.</p>
</dd>
</dl>

<dl>
<dt><code>__/j/db/add</code></dt>
<dd>
<p>
Adds documents to, or attributes to specific documents in, a collection.
<br>
The format of an Add, written in JSON, is as follows:
</p>

```json
{
    //Name of the collection to add to, will fail if omitted
    "collection": "some-collection",
    //List of full documents to add to the collection
    "documents": [
        {
            "_id": "some-new-id",
            "data": {
                "some-attribute": "some-value"
            }
        }
    ],
    //Map of new attributes that will be added to the documents they are listed under. Note that "_id" cannot be added as an attribute
    "values":{
        "some-document-id": {
            "some-new-attribute": "some-new-value"
        }
    }
}
```
<p>This request returns a list of elements added, and elements skipped.</p>
</dd>
</dl>

<dl>
    <dt><code>__/j/db/mod</code></dt>
    <dd>
    <p>
    Modifies attributes of a specific document within a collection.
    <br>
    The format of a Mod, written in JSON, is as follows:
    </p>

```json
{
    //Collection containing the document. Will fail if this collection does not exist
    "collection": "some-collection",
    //_id of Document to be modified. Will fail if this document does not exist
    "document": "some-document-id",
    //Map of attributes to be changed. Note that the _id can also be changed in this operation.
    "values":{
        "_id": "some-new-id",
        "some-attribute": "some-new-value"
    }
}
```
<p>This request returns a list of changes made, and changes skipped.</p>
</dd>
</dl>

<dl>
    <dt><code>__/j/db/delete</code></dt>
    <dd>
    <p>
    Deletes attributes of a specific document within a collection.
    <br>
    The format of a Delete, written in Json, is very similar to a Mod, as follows:
    </p>

```json
{
    //Collection containing the document. Will fail if this collection does not exist
    "collection": "some-collection",
    //_id of Document to be modified. Will fail if this document does not exist
    "document": "some-document-id",
    //List of attributes to be deleted. Note that _id cannot be deleted in this operation.
    "values":[
        "some-attribute",
        "another-attribute"
    ]
}
```

<p>This request returns a list of attributes deleted, and attributes skipped.</p>
</dd>
</dl>

<dl>
    <dt><code>__/j/db/purge</code></dt>
    <dd>
    <p>
    Deletes documents from a list of, or all, collections, based on query rules.
    <br>
    See <code>__/j/db/query</code> above, as this request has an identical structure.
    </p>
    <p>This request returns the list of document _ids deleted.</p>
    </dd>
</dl>
<br>
<h3>Skips</h3>
<p>
When performing operations such as Mod or Add, the request body contains rules for finding the targets of the operations. However, these may not always match up with the data, for example, trying to add a document whose _id already exists within that collection. In cases such as these, this part of the operation is skipped, and message is relayed back to the user through the HTTP Response data.
</p>
<br><br>
<h2>Services</h2>
<p>
Services are small python scripts that can be assigned to run on specific endpoints.
<br>This feature isn't implemented yet :(
</p>

<h2 align="center">License and Copyright Notice</h2>
<p align="center">
    Copyright (c) 2024, alchemicode. All Rights Reserved. 
    Permission to modify and redistribute is granted under the terms of the Apache 2.0.
</p>
