/* global $, Highcharts, _, Blob, URL, console */

(function(){
	"use strict";

	$(document).ready(function() {
		fetchDataAndRender();
		$('.viewoption').change(renderChart);
		$('.dataoption').change(fetchDataAndRender);
		$('#sharedTooltip').change(toggleSharedTooltip);
		$('#chartHeightContainer').change(setChartHeight);
		$('.dateContainer > div > span').click(setDates);
		$('#downloadData').click(downloadData);
	});

	function fetchDataAndRender(){
		$.ajax('/data/', {
			data: {from: $('#from-date').val().trim(), to:$('#to-date').val().trim()},
			error: function(err){console.log('something went wrong', err);},
			success: function(data){
				var browsers = $('#browsers'),
					limit = $('#browser-limit');
				window.chartData = data;
				window.dateSums = _(window.chartData).chain().pluck('data').flatten()
					.reduce(function(sums, p){
						sums[p.x] = p.y +(sums[p.x] || 0);
						return sums;
					}, {}).value();
				limit.empty();
				browsers.html($('<option>'))
				_(data).pluck('name').sort().forEach(function(name, idx){
					limit.append($('<option>').text(idx).attr('selected',idx===8));
					browsers.append($('<option>').text(name));
				});
				renderChart();
		    }
		});
	}

	function renderChart(){
		var obj = createHighchartsSettings(),
			chart = getChart();
		obj.series = getChartData();
		if(chart)
			chart.destroy();
		$('#chart-container').highcharts(obj);
	}

	function createHighchartsSettings() {
		return {
		    chart: {
		        renderTo: 'chart-container',
		        type: $('#chart-type').val().trim(),
		        zoomType: 'x',
		        height: $('#chartHeightContainer input[type=number]').val().trim()
		    },
		    title: {
				text: null
		    },
		    xAxis: {
		        type: 'datetime'
		    },
		    yAxis: {
		        title: {
		            text: 'Unique visitors'
		        },
		        min: 0
		    },
		    tooltip: {
				shared:true,
				useHTML:true,
				pointFormat: '<tr><td style="color:{series.color}">{series.name}:</td><td style="font-weight:bold;text-align:right;">{point.key}</td><td style="font-weight:bold;text-align:right;">({point.percentage}%%)</td></tr>',
				footerFormat: '<tr><td>Total visitors:</td><td style="font-weight:bold;text-align:right;">{total}</td></tr></table>',
				formatter: function(tooltip) {
					var points = this.points || [this],
						sum = window.dateSums[points[0].x],
						headerFormat = '{key:' + tooltip.options.dateTimeLabelFormats.week + '}<table>';
					return Highcharts.format(headerFormat, points[0]) +
						_(points).map(function(point){
							point.percentage= (point.y/sum < 0.1 ? '0' : '') + (100*point.y/sum).toFixed(2);
							point.key = Highcharts.numberFormat(point.y, 0, ',', ' ');
							return Highcharts.format(tooltip.options.pointFormat, {series: point.series, point: point});
						}).join('') + Highcharts.format(tooltip.options.footerFormat, {total:Highcharts.numberFormat(sum, 0, ',', ' ')});
				}
			},
		    plotOptions: {
				series: {
					stacking: $('#stacking').val().trim()
				}
		    },
		    credits: {
				enabled:false
		    }
		};
	}

	function getChartData(){
		var data = $.extend(true, [], window.chartData),
			browsers = $('#browsers').val(),
			limit = $('#browser-limit').val().trim();
		if(browsers && browsers[0]){
			data = _(data).filter(function(browser){
				return browsers.indexOf(browser.name) !== -1;
			});
		}
		if(limit)
			data = _(data).first(limit);
		return data;
	}

	function downloadData(){
		var a = document.createElement('a');
		a.href = URL.createObjectURL(
			new Blob(
				[JSON.stringify(
					window.chartData, function(key, value){
						return key==='x' ? new Date(value).toDateString() : value;
					}, '  ')
				],
				{"type" : "application\/octet-stream"}
			));
		a.download = "Data.txt";
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
	}

	function setDates(event){
		var period = event.target.dataset.period,
			date = new Date();
		if(period === 'always'){
			$('#from-date').val('');
			$('#to-date').val('');
			fetchDataAndRender();
			return;
		}
		date.setMinutes(date.getMinutes() - date.getTimezoneOffset());
		$('#to-date').val(date.toISOString().substring(0,10));
		date=new Date();
		switch(period){
		case 'month':
			date.setMonth(date.getMonth()-1);
			break;
		case '3month':
			date.setMonth(date.getMonth()-3);
			break;
		case '6month':
			date.setMonth(date.getMonth()-6);
			break;
		case 'year':
			date.setYear(1900+date.getYear()-1);
			break;
		case '3year':
			date.setYear(1900+date.getYear()-3);
			break;
		}
		date.setMinutes(date.getMinutes() - date.getTimezoneOffset());
		$('#from-date').val(date.toISOString().substring(0,10));
		fetchDataAndRender();
	}

	function setChartHeight(event){
		var el = event.target,
			type = el.type,
			value = el.value,
			chart = getChart();
		$('#chartHeightContainer input[type='+(type==='range' ? 'number' : 'range')+']').val(value);
		window.clearTimeout(window.timeoutId);
		window.timeoutId = window.setTimeout(function(){
			chart.setSize(chart.chartWidth, value);
		}, 250);
	}

	function toggleSharedTooltip(event){
		var checked = event.target.checked,
			chart = getChart();
		chart.options.tooltip.shared=checked;
		chart.tooltip = new Highcharts.Tooltip(chart, chart.options.tooltip);
	}

	function getChart(){
		return _.find(Highcharts.charts, function(c){return !!c;});
	}
})();