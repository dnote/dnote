/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import React from 'react';
import * as d3 from 'd3';
import _ from 'lodash';
import moment from 'moment';
import tippy from 'tippy.js';

import styles from './LessonHeatmap.module.scss';

function getSystemDate(systemUtcOffsetSeconds) {
  const date = new Date();
  const localUtcOffsetMinutes = 0 - date.getTimezoneOffset();
  const systemUtcOffsetMinutes = systemUtcOffsetSeconds / 60;
  date.setMinutes(
    date.getMinutes() - localUtcOffsetMinutes + systemUtcOffsetMinutes
  );
  return date;
}

function getDayDifference(a, b) {
  const millisecondsPerDay = 1000 * 60 * 60 * 24;
  const date1 = Date.UTC(a.getFullYear(), a.getMonth(), a.getDate());
  const date2 = Date.UTC(b.getFullYear(), b.getMonth(), b.getDate());

  return Math.floor((date2 - date1) / millisecondsPerDay);
}

// Given a date object returns the day of the week in English
export const getDayName = date =>
  [
    'Sunday',
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday',
    'Friday',
    'Saturday'
  ][date.getDay()];

function formatTooltipText({ date, count }) {
  const dateObject = new Date(date);
  const dateDayName = getDayName(dateObject);
  const dateText = moment(dateObject).format('MMM D YYYY');

  let contribText = 'No learning';
  if (count > 0) {
    contribText = `${count} learning${count > 1 ? 's' : ''}`;
  }
  return `<div class="learning-count">${contribText}</div><div class="learning-date">${dateDayName} ${dateText}</div>`;
}

export default class LessonHeatmap extends React.Component {
  constructor(props) {
    super(props);

    this.daySpace = 1;
    this.daySize = 12;

    this.daySizeWithSpace = this.daySize + this.daySpace * 2;
    this.monthNames = [
      'Jan',
      'Feb',
      'Mar',
      'Apr',
      'May',
      'Jun',
      'Jul',
      'Aug',
      'Sep',
      'Oct',
      'Nov',
      'Dec'
    ];
    this.months = [];
    this.group = 0;

    // Init color functions
    this.colorKey = this.initColorKey();
    this.color = this.initColor();

    // Loop through the timestamps to create a group of objects
    // The group of objects will be grouped based on the day of the week they are
    this.timestampsTmp = [];

    const today = getSystemDate(0);
    today.setHours(0, 0, 0, 0, 0);

    const oneYearAgo = new Date(today);
    oneYearAgo.setFullYear(today.getFullYear() - 1);

    const days = getDayDifference(oneYearAgo, today);

    for (let i = 0; i <= days; i += 1) {
      const date = new Date(oneYearAgo);
      date.setDate(date.getDate() + i);

      const day = date.getDay();
      const count = props.timestamps[moment(date).format('YYYY-M-D')] || 0;

      // Create a new group array if this is the first day of the week
      // or if is first object
      if ((day === 0 && i !== 0) || i === 0) {
        this.timestampsTmp.push([]);
        this.group += 1;
      }

      // Push to the inner array the values that will be used to render map
      const innerArray = this.timestampsTmp[this.group - 1];
      innerArray.push({ count, date, day });
    }
  }

  componentDidMount() {
    this.svg = this.renderSvg();
    this.renderDays();
    this.renderDayTitles();
    this.renderMonths();
    this.renderLegends();

    tippy('.js-tooltip', {
      arrow: true,
      theme: 'dnote'
    });
  }

  // Add extra padding for the last month label if it is also the last column
  getExtraWidthPadding = group => {
    let extraWidthPadding = 0;
    const lastColMonth = this.timestampsTmp[group - 1][0].date.getMonth();
    const secondLastColMonth = this.timestampsTmp[group - 2][0].date.getMonth();

    if (lastColMonth !== secondLastColMonth) {
      extraWidthPadding = 9;
    }

    return extraWidthPadding;
  };

  initColor = () => {
    const colorRange = [
      '#ededed',
      this.colorKey(0),
      this.colorKey(1),
      this.colorKey(2),
      this.colorKey(3)
    ];

    return d3
      .scaleThreshold()
      .domain([0, 2, 4, 6])
      .range(colorRange);
  };

  initColorKey = () =>
    d3
      .scaleLinear()
      .range(['#aecfe5', '#06435A'])
      .domain([0, 3]);

  renderSvg = () => {
    const width =
      (this.group + 1) * this.daySizeWithSpace +
      this.getExtraWidthPadding(this.group);

    return d3
      .select('#lesson-heatmap')
      .append('svg')
      .attr('width', width)
      .attr('height', 167)
      .attr('class', 'lesson-calendar');
  };

  renderDays = () => {
    this.svg
      .selectAll('g')
      .data(this.timestampsTmp)
      .enter()
      .append('g')
      .attr('transform', (group, i) => {
        _.each(group, (stamp, a) => {
          if (a === 0 && stamp.day === 0) {
            const month = stamp.date.getMonth();
            const x = this.daySizeWithSpace * i + 1 + this.daySizeWithSpace;
            const lastMonth = _.last(this.months);
            if (
              lastMonth == null ||
              (month !== lastMonth.month &&
                x - this.daySizeWithSpace !== lastMonth.x)
            ) {
              this.months.push({ month, x });
            }
          }
        });
        return `translate(${this.daySizeWithSpace * i +
          1 +
          this.daySizeWithSpace}, 18)`;
      })
      .selectAll('rect')
      .data(stamp => stamp)
      .enter()
      .append('rect')
      .attr('x', '0')
      .attr('y', stamp => this.daySizeWithSpace * stamp.day)
      .attr('width', this.daySize)
      .attr('height', this.daySize)
      .attr('fill', stamp =>
        stamp.count !== 0 ? this.color(Math.min(stamp.count, 40)) : '#ededed'
      )
      .attr('title', stamp => formatTooltipText(stamp))
      .attr('class', 'js-tooltip')
      .on('click', this.clickDay);
  };

  renderDayTitles = () => {
    const days = [
      {
        text: 'M',
        y: 29 + this.daySizeWithSpace * 1
      },
      {
        text: 'W',
        y: 29 + this.daySizeWithSpace * 3
      },
      {
        text: 'F',
        y: 29 + this.daySizeWithSpace * 5
      }
    ];

    this.svg
      .append('g')
      .selectAll('text')
      .data(days)
      .enter()
      .append('text')
      .attr('text-anchor', 'middle')
      .attr('x', 8)
      .attr('y', day => day.y)
      .text(day => day.text)
      .attr('class', 'user-contrib-text');
  };

  renderMonths = () => {
    this.svg
      .append('g')
      .attr('direction', 'ltr')
      .selectAll('text')
      .data(this.months)
      .enter()
      .append('text')
      .attr('x', date => date.x)
      .attr('y', 10)
      .attr('class', 'user-contrib-text')
      .text(date => this.monthNames[date.month]);
  };

  renderLegends = () => {
    const keyValues = [
      '<div class="learning-count">No learning</div>',
      '<div class="learning-count">1 learning</div>',
      '<div class="learning-count">2-3 learnings</div>',
      '<div class="learning-count">4-5 learnings</div>',
      '<div class="learning-count">5+ learnings</div>'
    ];

    const keyColors = [
      '#ededed',
      this.colorKey(0),
      this.colorKey(1),
      this.colorKey(2),
      this.colorKey(3)
    ];

    this.svg
      .append('g')
      .attr('transform', `translate(18, ${this.daySizeWithSpace * 8 + 16})`)
      .selectAll('rect')
      .data(keyColors)
      .enter()
      .append('rect')
      .attr('width', this.daySize)
      .attr('height', this.daySize)
      .attr('x', (color, i) => this.daySizeWithSpace * i)
      .attr('y', 0)
      .attr('fill', color => color)
      .attr('class', 'js-tooltip')
      .attr('title', (color, i) => keyValues[i])
      .attr('data-container', 'body');
  };

  render() {
    return (
      <div className="calendar">
        <div id="lesson-heatmap" className={styles.heatmap} />
      </div>
    );
  }
}
