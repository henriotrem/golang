var width = 600;
var height = 5000;
var barHeight = 20;
var barWidth = width * .8;

var intendedIdx = 0,
    duration = 350,
    intendedRoot;

var intendedTreeMap = d3.tree().size([width, height])
    .nodeSize([0, 30]);

var intendedSvg = d3.select('#indexes').append('svg')
    .attr('width', width)
    .attr('height', height)
    .append('g')
    .attr('transform', 'translate(10,20)');;

function updateIntendedData() {

    intendedIdx = 0

    intendedRoot = intendedTreeMap(d3.hierarchy(intendedTreeData));

    intendedRoot.each(function (d, idx) {
        d.name = d.data.name; //transferring name to a name variable
        d.id = idx; 
    });

    intendedRoot.x0 = intendedRoot.x;
    intendedRoot.y0 = intendedRoot.y

    //root.children.forEach(intendedCollapse);
    updateIntendedTree(intendedRoot);

}

function updateIntendedTree(source) {

    // Compute the new tree layout.
    let nodes = intendedTreeMap(intendedRoot)
    let nodesSort = [];
    nodes.eachBefore(function (n) {
        nodesSort.push(n);
    });
    height = Math.max(600, nodesSort.length);
    let links = nodesSort.slice(1);
    // Compute the "layout".
    nodesSort.forEach(function(n, i) {
        n.x = i * barHeight;
    });

    d3.select('svg').transition()
        .duration(duration)
        .attr("height", height);

    // Update the nodes…
    let node = intendedSvg.selectAll('g.node')
        .data(nodesSort, function (d) {
            return d.id || (d.id = ++intendedIdx);
        });

    // Enter any new nodes at the parent's previous position.
    let nodeEnter = node.enter().append('g')
        .attr('class', 'node')
        .attr('transform', function () {
            return 'translate(' + source.y0 + ',' + source.x0 + ')';
        })
        .on('click', intendedClick);

    nodeEnter.append('circle')
        .attr('r', 1e-6)
        .style('fill', function (d) {
            return d._children ? 'lightsteelblue' : '#fff';
        });

    nodeEnter.append('text')
        .attr('x', function (d) {
            return d.children || d._children ? 10 : 10;
        })
        .attr('dy', '.35em')
        .attr('text-anchor', function (d) {
            return d.children || d._children ? 'start' : 'start';
        })
        .text(function (d) {
            if (d.data.name.length > 20) {
                return d.data.name.substring(0, 20) + '...';
            } else {
                return d.data.name;
            }
        })
        .style('fill-opacity', 1e-6);

    nodeEnter.append('svg:title').text(function (d) {
        return d.data.name;
    });

    // Transition nodes to their new position.
    let nodeUpdate = node.merge(nodeEnter)
        .transition()
        .duration(duration);

    nodeUpdate
        .attr('transform', function (d) {
            return 'translate(' + d.y + ',' + d.x + ')';
        });

    nodeUpdate.select('circle')
        .attr('r', 4.5)
        .style('fill', function (d) {
            return d._children ? 'lightsteelblue' : '#fff';
        });

    nodeUpdate.select('text')
        .style('fill-opacity', 1);

    // Transition exiting nodes to the parent's new position (and remove the nodes)
    let nodeExit = node.exit().transition()
        .duration(duration);

    nodeExit
        .attr('transform', function (d) {
            return 'translate(' + source.y + ',' + source.x + ')';
        })
        .remove();

    nodeExit.select('circle')
        .attr('r', 1e-6);

    nodeExit.select('text')
        .style('fill-opacity', 1e-6);

    // Update the links…
    let link = intendedSvg.selectAll('path.link')
        .data(links, function (d) {
            // return d.target.id;
            let id = d.id + '->' + d.parent.id;
            return id;
        });

    // Enter any new links at the parent's previous position.
    let linkEnter = link.enter().insert('path', 'g')
        .attr('class', 'link')
        .attr('d', function (d) {
            let o = {
                x: source.x0,
                y: source.y0,
                parent: {
                    x: source.x0,
                    y: source.y0
                }
            };
            return intendedConnector(o);
        });

    // Transition links to their new position.
    link.merge(linkEnter).transition()
        .duration(duration)
        .attr('d', intendedConnector);


    // // Transition exiting nodes to the parent's new position.
    link.exit().transition()
        .duration(duration)
        .attr('d', function (d) {
            let o = {
                x: source.x,
                y: source.y,
                parent: {
                    x: source.x,
                    y: source.y
                }
            };
            return intendedConnector(o);
        })
        .remove();

    // Stash the old positions for transition.
    nodesSort.forEach(function (d) {
        d.x0 = d.x;
        d.y0 = d.y;
    });
}

function intendedConnector(d) {
    return "M" + d.parent.y + "," + d.parent.x +
        "V" + d.x + "H" + d.y;
};

function intendedCollapse(d) {
    if (d.children) {
        d._children = d.children;
        d._children.forEach(intendedCollapse);
        d.children = null;
    }
};

function intendedClick(evt, d) {
    if (d.children) {
        d._children = d.children;
        d.children = null;
    } else {
        d.children = d._children;
        d._children = null;
    }
    updateIntendedTree(d);
};