var serverLogs = document.getElementById("server")
var inputLogs = document.getElementById("inputLogs");
var outputLogs = document.getElementById("outputLogs");
var selectList = document.getElementById("select");
var ws;

var radialTreeData =  {"name": "#", "count": 0}
updateRadialData()

var intendedTreeData =  {"name": "#", "count": 0}
updateIntendedData()

var network = []

var servers = {}
var trees = {}

var input = {}
var output = {}
var scanned = {}

var ip = ""

var createNetwork = function() {
  network = []
  for (let idx in servers) {
    var object = {name: idx, imports:[]}
    for (var distant in servers[idx].address_lookup) {
      if(distant != idx) {
        object.imports.push(distant)
      }
    }
    network.push(object)
  }
}

var scanTree = function(serverNode) {
  
  scanned[ip][serverNode.name] = servers[ip].address.leaves[serverNode.name] != null ? "local" : serverNode.children != null ? "aggregated" : "distant"

  if (serverNode.children) {
    for (var child of serverNode.children) {          
      scanTree(child)
    }
  }
}

var filterNodes = function(networkNode) {
  
  networkNode["type"] = scanned[ip][networkNode.name] != null ? scanned[ip][networkNode.name] : "none"

  if (networkNode.children) {
    for (var child of networkNode.children) {          
      filterNodes(child)
    }
  }
}

var addLeaf = function(node, parent, leaf) {
  
  if (node.name == parent) {
    if (node.children == null) {
      node.children = []
    }
    if (node.children.length < 16) {
      node.children.push({"name":leaf, "count":0})
    }
    return true
  }

  if (node.children) {
    for (var child of node.children) {
      if (addLeaf(child, parent, leaf)) {
        return true
      }
    }
  }

  return false
}

var updateRoot = function(lookup) {
  for (const [key, value] of Object.entries(lookup)) {
    for (const [subkey, subvalue] of Object.entries(value.leaves)) {
      addLeaf(radialTreeData, subkey.substring(0, subkey.length-1), subkey)
    }
  }
}

var printInput = function(message) {
  var d = document.createElement("div");
  d.innerHTML = message;
  inputLogs.insertBefore(d, inputLogs.firstChild)
  if (inputLogs.length > 40) {
    inputLogs.removeChild(inputLogs.lastElementChild);
  }
};

var printOutput = function(message) {
  var d = document.createElement("div");
  d.innerHTML = message;
  outputLogs.insertBefore(d, outputLogs.firstChild)
  if (outputLogs.length > 40) {
    outputLogs.removeChild(outputLogs.lastElementChild);
  }
};

var newServer = function(newIp) {
  var option = document.createElement("option");
  option.setAttribute("value", newIp);
  option.text = newIp;
  selectList.appendChild(option);
  selectList.value = newIp;

  output[newIp] = []
  input[newIp] = []
  scanned[newIp] = {}
}

var refreshLogs = function() {
  serverLogs.innerHTML =  JSON.stringify(servers[ip], undefined, 2);
  inputLogs.innerHTML = ""
  outputLogs.innerHTML = ""

  for (var log of input[ip]) {
    printInput(log)
  }
  for (var log of output[ip]) {
    printOutput(log)
  }
}

var refreshTree = function(ip) {
  if (trees[ip] != null) {
    console.log(trees[ip])
    intendedTreeData = trees[ip]
    updateIntendedData()
    scanTree(intendedTreeData)  
    filterNodes(radialTreeData);    
    updateRadialData()
  }
}

document.getElementById("select").onchange = function(evt) {
  ip = selectList.value
  refreshLogs()
  refreshTree(ip)
}

document.getElementById("open").onclick = function(evt) {
  if (ws) {
    return false;
  }
  var loc = window.location, new_uri;
  if (loc.protocol === "https:") {
    new_uri = "wss:";
  } else {
    new_uri = "ws:";
  }
  new_uri += "//" + loc.host;
  new_uri += loc.pathname + "ws";
  ws = new WebSocket(new_uri);
  ws.onopen = function(evt) {
    console.log("OPEN");
  }
  ws.onclose = function(evt) {
    console.log("CLOSE");
    ws = null;
  }
  ws.onmessage = function(evt) {

    var response = evt.data.split('\t')

    if(response[0] == "NETWORK") {

      var origin = response[1]
      var destination = response[2]
      var type = response[3]
      var content = response[4]

      if (type == "main.ServeAction") { 
        updateRoot(JSON.parse(content).address_lookup)
        createNetwork()
        updateHierarchicalBundle()
      }
      
      if (scanned[origin] == null) {
        ip = origin
        newServer(ip)
      }

      if (scanned[destination] == null) {
        ip = destination
        newServer(ip)
      }

      input[destination].push(origin + " " + type)
      if (input[destination].length > 40) {
        input[destination].shift()
      }

      output[origin].push(destination + " " + type)
      if (output[origin].length > 40) {
        output[origin].shift()
      }
      
      if (origin == ip || destination == ip) {
        refreshLogs()
      }

    } else if (response[0] == "STATE") {

      var origin = response[1]
      var server = JSON.parse(response[2])
      var tree = JSON.parse(response[3])
      var bool = false

      if (trees[origin] == null) {
        bool = true
      }

      servers[origin] = server
      trees[origin] = tree
      
      if (bool) {
        refreshTree(origin)
      }
    }
  }
  ws.onerror = function(evt) {
    console.log("ERROR: " + evt.data);
  }
  return false;
};

document.getElementById("init").onclick = function(evt) {
  if (!ws) {
    return false;
  }
  console.log("SEND: INIT");
  ws.send("INIT");
  return false;
};

document.getElementById("new").onclick = function(evt) {
  if (!ws) {
    return false;
  }
  ws.send("NEW");
  return false;
};

document.getElementById("quit").onclick = function(evt) {
  evt.preventDefault();
  if (!ws) {
    return false;
  }
  console.log("SEND: QUIT");
  ws.send("QUIT");
  return false;
};

document.getElementById("close").onclick = function(evt) {
  if (!ws) {
    return false;
  }
  ws.close();
  return false;
};