var diameter = 600;
    width = diameter,
    height = diameter;
    
var i = 0,
    duration = 350,
    root;

var radialTreeMap = d3.tree()
    .size([360, diameter / 2 - 80])
    .separation(function(a, b) { return (a.parent == b.parent ? 1 : 10) / a.depth; });

var radialTreeSvg = d3.select("#radial").append("svg")
    .attr("width", width )
    .attr("height", height )
    .append("g")
    .attr("transform", "translate(" + diameter / 2 + "," + diameter / 2 + ")");

function updateRadialData() {

    i = 0

    root = d3.hierarchy(radialTreeData, function(d) {
        return d.children;
    });

    root.each(function (d) {
        d.name = d.data.name; //transferring name to a name variable
        d.type = d.data.type;
        d.count = d.data.count;
        d.id = i; //Assigning numerical Ids
        i += i;
    });

    root.x0 = height / 2;
    root.y0 = 0;

    updateRadialTree(root)
}

function updateRadialTree(source) {

  var nodes = radialTreeMap(root).descendants();
  var links = nodes.slice(1);

  // Normalize for fixed-depth.
  nodes.forEach(function(d) { d.y = d.depth * 80; });

  // Update the nodes…
  var node = radialTreeSvg.selectAll("g.radialNode")
      .data(nodes, function(d) { return d.id || (d.id = ++i); });

  // Enter any new nodes at the parent's previous position.
  var nodeEnter = node.enter().append("g")
      .attr("class", "radialNode")      
      .attr('d', function(d){
        var o = {x: source.x0, y: source.y0}
        return diagonal(o, o)
      })
      .on("click", click);

  nodeEnter.append("circle")
      .attr("r", 4)
      .style("fill", function(d) { return d._children ? "lightsteelblue" : "#fff"; })
      .style("stroke", function(d) {  return d.type == "local" ? "#C51162" : d.type == "aggregated" ? "#4DB6AC" : d.type == "distant" ? "#2196F3" : "black" ; });

  nodeEnter.append("text")
      .style("fill-opacity", 1)
      .attr("transform", function(d) { return d.x < 180 ? "translate(10)" : "rotate(180)translate(-" + (d.name.length*3+10)  + ")"; })
      .text(function(d) {  return d.name; });
      
  var nodeUpdate = nodeEnter.merge(node);

  // Transition nodes to their new position.
  nodeUpdate.transition()
    .duration(duration)
    .attr("transform", function(d) { return "rotate(" + (d.x - 90) + ")translate(" + d.y + ")"; });

  nodeUpdate.select("circle")
    .attr("r", 4)
    .style("fill", function(d) { return d._children ? "lightsteelblue" : "#fff"; })
    .style("stroke", function(d) { return d.type == "local" ? "#C51162" : d.type == "aggregated" ? "#4DB6AC" : d.type == "distant" ? "#2196F3" : "black" ; });

  nodeUpdate.select("text")
    .style("fill-opacity", 1)
    .attr("transform", function(d) { return d.x < 180 ? "translate(10)" : "rotate(180)translate(-" + (d.name.length*3+20)  + ")"; });

  // TODO: appropriate transform
  var nodeExit = node.exit().transition()
    .duration(duration)
    .remove();

  nodeExit.select("circle")
      .attr("r", 1e-6);

  nodeExit.select("text")
      .style("fill-opacity", 1e-6);

  // Update the links…
  var link = radialTreeSvg.selectAll("path.radialLink")
      .data(links, function(d) { return d.id; });

  // Enter any new links at the parent's previous position.
  var linkEnter = link.enter().insert("path", "g")
      .attr("class", "radialLink")
      .attr('d', function(d){
        var o = {x: source.x0, y: source.y0}
        return diagonal(o, o)
      });

  var linkUpdate = linkEnter.merge(link);

  // Transition links to their new position.
  linkUpdate.transition()
    .duration(duration)
    .attr("d",  function(d) { return diagonal(d, d.parent)});

  // Transition exiting nodes to the parent's new position.
  var linkExit = link.exit().transition()
    .duration(duration)
    .attr('d', function(d) {
        var o = {x: source.x, y: source.y}
        return diagonal(o, o)
    })
    .remove();

  // Stash the old positions for transition.
  nodes.forEach(function(d) {
    d.x0 = d.x;
    d.y0 = d.y;
  });
}

// Creates a curved (diagonal) path from parent to the child nodes
function diagonal(source, destination) {

    path = "M" + project(source.x, source.y)
            + "C" + project(source.x, (source.y + destination.y) / 2)
            + " " + project(destination.x, (source.y + destination.y) / 2)
            + " " + project(destination.x, destination.y)

    return path
}


function project(x, y) {
    var angle = (x - 90) / 180 * Math.PI, radius = y;
    return [radius * Math.cos(angle), radius * Math.sin(angle)];
}

// Toggle children on click.
function click(evt, d) {
  if (d.children) {
    d._children = d.children;
    d.children = null;
  } else {
    d.children = d._children;
    d._children = null;
  }
  
  updateRadialTree(d);
}

// Collapse nodes
function collapse(d) {
  if (d.children) {
      d._children = d.children;
      d._children.forEach(collapse);
      d.children = null;
    }
}