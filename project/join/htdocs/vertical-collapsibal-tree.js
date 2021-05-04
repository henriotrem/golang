var width = 600 ,
  height = 1000;

var duration = 350,
  verticalRoot;

var verticalTreeMap = d3.tree()
  .nodeSize([35, 150])

var verticalSvg = d3.select("#indexes").append("svg")
  .attr("width", width)
  .attr("height", height)
  .append("g")
  .attr("transform",
        "translate(300,50)");

function updateVerticalData() {
  
  verticalRoot = d3.hierarchy(verticalTreeData); 
  updateVerticalTree(verticalRoot);
}

function updateVerticalTree(source) {

  var nodes = verticalTreeMap(verticalRoot).descendants()

  nodes.forEach(function(node, idx) {
        node.id = idx; 
        if (node.parent != null && node.parent.parent != null) {
            node.x = node.parent.parent.children[(node.id-1) % node.parent.children.length].x; 
        }
    });

  var links = nodes.slice(1);
  
  var node = verticalSvg.selectAll("g.node")
      .data(nodes);

  let nodeEnter = node.enter().append('g')
      .attr('class', 'node')
      .attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; })
      .on('click', verticalClick);

  nodeEnter.append('circle')
      .attr("r", 4.5)
      .style("stroke", function(d) { return d._children || d.children ? "#4DB6AC" : "#2196F3" ; });

  nodeEnter.append('text')
      .attr("dy", ".35em")
      .attr("y", function(d) { return  d._children || d.children ? -20 : 20; })
      .style("text-anchor", "middle")
      .text(function(d) { return d.data.name; })
      .style('fill-opacity', 1e-6);

  let nodeUpdate = node.merge(nodeEnter)
      .transition()
      .duration(duration);
      
  nodeUpdate.select('circle')
    .attr("r", 4.5)
    .style("stroke", function(d) { return d._children || d.children ? "#4DB6AC" : "#2196F3" ; });
  
  nodeUpdate.select('text')
      .attr("dy", ".35em")
      .attr("y", function(d) { return  d._children || d.children ? -20 : 20; })
      .text(function(d) { return d.data.name; })
      .style('fill-opacity', 1);

  node.exit().remove();

  let link = verticalSvg.selectAll('path.link')
      .data(links);

  let linkEnter =  link.enter().insert('path', 'g')
      .attr('class', 'link');

  link.merge(linkEnter)
      .attr('d', verticalConnector);

  link.exit().remove();
}

function verticalConnectorStraight(d) {
  return "M" + (d.x) + "," + d.y
       + " " + (d.parent.x) + "," + d.parent.y;
};

function verticalConnector(d) {
  return "M" + d.x + "," + d.y
      + "C" + d.x + "," + (d.y + d.parent.y) / 2
      + " " + d.parent.x + "," +  (d.y + d.parent.y) / 2
      + " " + d.parent.x + "," + d.parent.y;
};

function verticalClick(evt, d) {
  if (d.children) {
      d._children = d.children;
      d.children = null;
  } else {
      d.children = d._children;
      d._children = null;
  }
  updateVerticalTree(d);
};