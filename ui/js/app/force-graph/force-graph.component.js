angular.
module('covalence').
component('forceGraph', {
    templateUrl: 'js/app/force-graph/force-graph.template.html',
    controller: function ForceGraphController($scope, $http, $timeout) {
        $scope.deploymentFilters = {}
        $scope.vmsFilters = {}
        var w = 680,
            h = 700,
            rx = w / 2,
            ry = h / 2,
            m0,
            rotate = 0;
        // Chrome 15 bug: <http://code.google.com/p/chromium/issues/detail?id=98951>
        var div = d3.select("div.force").insert("div", "h2")
            .style("top", "0px")
            .style("left", "0px")
            .style("width", w + "px")
            .style("height", h + "px")
            .style("position", "relative")
            .style("-webkit-backface-visibility", "hidden");

        var svg = div.append("svg:svg")
            .attr("width", w)
            .attr("height", h)

        var color = d3.scaleOrdinal(d3.schemeCategory20);

        var simulation = d3.forceSimulation()
            .force("link", d3.forceLink().id(function(d) {
                return d.id;
            }))
            .force("charge", d3.forceManyBody())
            .force("center", d3.forceCenter(w / 2, h / 2));
        $scope.filterDeployment = function(key) {
            $scope.deployments[key].vms.forEach(function(d) {
                $scope.vmsFilters[d.job_name + "/" + d.index] = $scope.deploymentFilters[key]
            })
            drawGraph($scope.connections)
        }
        $scope.filterVMs = function(key) {
            drawGraph($scope.connections)
        }
        $scope.showVms = function() {
            timeOffset = Math.floor(Date.now() / 1000) - 35;
            $http({
                method: 'GET',
                url: '/vms',
                params: {
                    "after": timeOffset
                }
            }).then(function successCallback(response) {
                    $scope.deployments = {}
                    var id = 1
                    $scope.vms = response.data
                    $scope.vms_map = {}
                    $scope.vms.forEach(function(d) {
                        $scope.vms_map[d.ip] = d.job_name + "/" + d.index
                    });
                    for (var i = 0; i < $scope.vms.length; i++) {
                        var name = $scope.vms[i].deployment_name
                        if ($scope.deployments[name] == undefined) {
                            $scope.deployments[name] = {
                                id: id,
                                vms: []
                            }
                            id++;
                        }
                        $scope.deployments[name].vms.push($scope.vms[i])
                    }
                },
                function errorCallback(response) {
                    // called asynchronously if an error occurs
                    // or server returns response with an error status.
                });

            $http({
                method: 'GET',
                url: '/connections',
                params: {
                    "after": timeOffset
                }
            }).then(function successCallback(response) {
                $scope.connections = response.data
                drawGraph($scope.connections)

            }, function errorCallback(response) {
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };

        function drawGraph(connections) {
            var nodes = {};
            var links = [];
            connections.filter(userFilter).forEach(function(value) {
                source_name = getName(value.source.ip)
                destination_name = getName(value.destination.ip)
                if (nodes[source_name] == undefined) {
                    node = {
                        id: source_name,
                        group: getDeployment(value.source.deployment)
                    }
                    nodes[source_name] = node
                }
                if (nodes[destination_name] == undefined) {
                    node = {
                        id: destination_name,
                        group: getDeployment(value.source.deployment)
                    }
                    nodes[destination_name] = node
                }
                link = {
                    source: getName(value.source.ip),
                    target: getName(value.destination.ip),
                    value: 20
                }
                links.push(link)
            });
            var nodes_array = []
            for (var key in nodes) {
                nodes_array.push(nodes[key])
            }
            d3.selectAll('g').remove()
            simulation.restart()
            var link = svg.append("g")
                .attr("class", "links")
                .selectAll("line")
                .data(links)
                .enter().append("line")

            var node = svg.append("g")
                .attr("class", "nodes")
                .selectAll("circle")
                .data(nodes_array)
                .enter().append("circle")
                .attr("r", 5)
                .attr("fill", function(d) {
                    return color(d.group);
                })
                .call(d3.drag()
                    .on("start", dragstarted)
                    .on("drag", dragged)
                    .on("end", dragended));

            node.append("title")
                .text(function(d) {
                    return d.id;
                });

            // add the text
            node.append("text")
                .attr("x", 12)
                .attr("dy", ".35em")
                .text(function(d) { return d.id; });
            simulation
                .nodes(nodes_array)
                .on("tick", ticked);

            simulation.force("link")
                .links(links);

            function ticked() {
                link
                    .attr("x1", function(d) {
                        return d.source.x;
                    })
                    .attr("y1", function(d) {
                        return d.source.y;
                    })
                    .attr("x2", function(d) {
                        return d.target.x;
                    })
                    .attr("y2", function(d) {
                        return d.target.y;
                    });

                node
                    .attr("cx", function(d) {
                        return d.x;
                    })
                    .attr("cy", function(d) {
                        return d.y;
                    });
            }
        }
        $scope.showVms();
        $scope.intervalFunction = function() {
            $timeout(function() {
                $scope.showVms();
                $scope.intervalFunction();
            }, 60000);

        };

        function userFilter(value) {
            if ($scope.vmsFilters[value.source.job + "/" + value.source.index]) {
                return value
            }
        }

        function getName(ip) {
            if ($scope.vms_map[ip]) {
                return $scope.vms_map[ip]
            }
            return ip
        }

        function getDeployment(deployment) {
            if ($scope.deployments[deployment]) {
                return $scope.deployments[deployment].id
            }
            return 20
        }



        function dragstarted(d) {
            if (!d3.event.active) simulation.alphaTarget(0.3).restart();
            d.fx = d.x;
            d.fy = d.y;
        }

        function dragged(d) {
            d.fx = d3.event.x;
            d.fy = d3.event.y;
        }

        function dragended(d) {
            if (!d3.event.active) simulation.alphaTarget(0);
            d.fx = null;
            d.fy = null;
        }

        $scope.intervalFunction();
    }
});