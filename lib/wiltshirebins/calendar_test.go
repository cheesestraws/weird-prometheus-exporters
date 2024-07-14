package wiltshirebins

var testCalendar = `
<div class="collection-calendar">
    <header class="collection-calendar-header">
        <form method="post">
            <input type="hidden" name="Postcode" value="SP1 3YG" />
            <input type="hidden" name="Uprn" value="100121037470" />
            <input type="hidden" id="PreviousMonth" name="Month" value="6" />
            <input type="hidden" id="PreviousYear" name="Year" value="2024" />
            <button type="submit" class="btn btn-primary rc-previous" value=""><i class="fa fa-chevron-left"></i> June 2024</button>
        </form>
        <form method="post">
            <input type="hidden" name="Postcode" value="SP1 3YG" />
            <input type="hidden" name="Uprn" value="100121037470" />
            <input type="hidden" id="NextMonth" name="Month" value="8" />
            <input type="hidden" id="NextYear" name="Year" value="2024" />
            <button type="submit" class="btn btn-primary rc-next" value="">August 2024 <i class="fa fa-chevron-right"></i></button>
        </form>
        <span class="rc-current-month">July 2024</span>
    </header>
    <div class="collection-calendar-content">
        <div id="calendar" class="cal-context" style="width: 100%;">
            <div class="cal-row-fluid cal-row-head">
                <div class="cal-cell1">Sun</div>
                <div class="cal-cell1">Mon</div>
                <div class="cal-cell1">Tues</div>
                <div class="cal-cell1">Wed</div>
                <div class="cal-cell1">Thurs</div>
                <div class="cal-cell1">Fri</div>
                <div class="cal-cell1">Sat</div>
            </div>

            <div class="cal-month-box">

                                <div class="cal-row-fluid cal-before-eventlist">
                                    <div class="cal-cell1 cal-cell cal-cell-outmonth" data-cal-row="-day0">
                                        <div class="cal-month-day cal-day-outmonth cal-day-weekend cal-month-first-row"><span class="pull-right" data-cal-date="" data-cal-view="day" data-toggle="tooltip" title="" data-original-title=""></span></div>
                                    </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day1">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-01T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">1</span><span class="day-name">MON</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day2">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-02T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">2</span><span class="day-name">TUE</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day3">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-03T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">3</span><span class="day-name">WED</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day4">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-04T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">4</span><span class="day-name">THU</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell cal-cell-active" data-cal-row="-day5">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-05T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">5</span><span class="day-name">FRI</span></div>
                                                    <div class="events-list">
                                                                <div class="rc-event-container"><a data-event-id="pod" class="event service-pod" data-toggle="tooltip" title="" data-original-datetext="Friday 5 July, 2024" data-original-title="Mixed dry recycling (blue lidded bin) and glass (black box or basket)" data-original-warning=""></a><span>Mixed dry recycling (blue lidded bin) and glass (black box or basket)</span></div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day6">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-06T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">6</span><span class="day-name">SAT</span></div>
                                                </div>
                                            </div>
                                        </div>
                                </div>
                                <div class="cal-row-fluid cal-before-eventlist">
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day0">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-07T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">7</span><span class="day-name">SUN</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day1">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-08T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">8</span><span class="day-name">MON</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day2">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-09T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">9</span><span class="day-name">TUE</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day3">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-10T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">10</span><span class="day-name">WED</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day4">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-11T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">11</span><span class="day-name">THU</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell cal-cell-active" data-cal-row="-day5">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-12T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">12</span><span class="day-name">FRI</span></div>
                                                    <div class="events-list">
                                                                <div class="rc-event-container"><a data-event-id="res" class="event service-res" data-toggle="tooltip" title="" data-original-datetext="Friday 12 July, 2024" data-original-title="Household waste" data-original-warning=""></a><span>Household waste</span></div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day6">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-13T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">13</span><span class="day-name">SAT</span></div>
                                                </div>
                                            </div>
                                        </div>
                                </div>
                                <div class="cal-row-fluid cal-before-eventlist">
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day0">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-14T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">14</span><span class="day-name">SUN</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day1">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-15T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">15</span><span class="day-name">MON</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day2">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-16T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">16</span><span class="day-name">TUE</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day3">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-17T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">17</span><span class="day-name">WED</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day4">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-18T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">18</span><span class="day-name">THU</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell cal-cell-active" data-cal-row="-day5">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-19T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">19</span><span class="day-name">FRI</span></div>
                                                    <div class="events-list">
                                                                <div class="rc-event-container"><a data-event-id="pod" class="event service-pod" data-toggle="tooltip" title="" data-original-datetext="Friday 19 July, 2024" data-original-title="Mixed dry recycling (blue lidded bin) and glass (black box or basket)" data-original-warning=""></a><span>Mixed dry recycling (blue lidded bin) and glass (black box or basket)</span></div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day6">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-20T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">20</span><span class="day-name">SAT</span></div>
                                                </div>
                                            </div>
                                        </div>
                                </div>
                                <div class="cal-row-fluid cal-before-eventlist">
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day0">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-21T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">21</span><span class="day-name">SUN</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day1">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-22T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">22</span><span class="day-name">MON</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day2">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-23T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">23</span><span class="day-name">TUE</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day3">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-24T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">24</span><span class="day-name">WED</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day4">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-25T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">25</span><span class="day-name">THU</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell cal-cell-active" data-cal-row="-day5">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-26T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">26</span><span class="day-name">FRI</span></div>
                                                    <div class="events-list">
                                                                <div class="rc-event-container"><a data-event-id="res" class="event service-res" data-toggle="tooltip" title="" data-original-datetext="Friday 26 July, 2024" data-original-title="Household waste" data-original-warning=""></a><span>Household waste</span></div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day6">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-27T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">27</span><span class="day-name">SAT</span></div>
                                                </div>
                                            </div>
                                        </div>
                                </div>
                                <div class="cal-row-fluid cal-before-eventlist">
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day0">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-28T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">28</span><span class="day-name">SUN</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day1">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-29T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">29</span><span class="day-name">MON</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day2">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-30T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">30</span><span class="day-name">TUE</span></div>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="cal-cell1 cal-cell" data-cal-row="-day3">
                                            <div class="cal-month-day cal-day-inmonth">
                                                <div class="cal-inner">
                                                    <div class="rc-day-of-week"><span class="day-no" data-cal-date="2024-07-31T00:00:00" data-cal-view="day" data-toggle="tooltip" title="" data-original-title="">31</span><span class="day-name">WED</span></div>
                                                </div>
                                            </div>
                                        </div>
                                    <div class="cal-cell1 cal-cell cal-cell-outmonth" data-cal-row="-day4">
                                        <div class="cal-month-day cal-day-outmonth cal-day-weekend cal-month-first-row"><span class="pull-right" data-cal-date="" data-cal-view="day" data-toggle="tooltip" title="" data-original-title=""></span></div>
                                    </div>
                                    <div class="cal-cell1 cal-cell cal-cell-outmonth" data-cal-row="-day5">
                                        <div class="cal-month-day cal-day-outmonth cal-day-weekend cal-month-first-row"><span class="pull-right" data-cal-date="" data-cal-view="day" data-toggle="tooltip" title="" data-original-title=""></span></div>
                                    </div>
                                    <div class="cal-cell1 cal-cell cal-cell-outmonth" data-cal-row="-day6">
                                        <div class="cal-month-day cal-day-outmonth cal-day-weekend cal-month-first-row"><span class="pull-right" data-cal-date="" data-cal-view="day" data-toggle="tooltip" title="" data-original-title=""></span></div>
                                    </div>
                                </div>


            </div>
            <div class="calendar-foot-section"></div>
            <div class="calendar-foot-panel panel-closed" style="max-height: 115px;">
                <div class="panel rc-calendar-day-panel">
                    <div class="panel-heading">
                        <header>
                            <h3><span id="rc-panel-date"> </span> <i class="fa fa-angle-double-up close-panel"></i></h3>
                        </header>
                    </div>
                    <div class="panel-body">
                        <ul class="list-unstyled" id="rc-panel-collection-list">
                            <li> </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
            <div class="rc-key">
                <div class="row">
                        <div class="col-sm-6">
                            <strong>
                                <span class="event service-pod"></span> Mixed dry recycling (blue lidded bin) and glass (black box or basket)
                            </strong>
                        </div>
                        <div class="col-sm-6">
                            <strong>
                                <span class="event service-res"></span> Household waste
                            </strong>
                        </div>
                </div>
                <p><a class="govuk-link govuk-link--no-visited-state" href="/wastecollectiondays/printablecalendar/100121037470"> Download a printable collections calendar</a></p>
            </div>

</div>
`
