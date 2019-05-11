package main

import (
	"net"
	"log"
	"bufio"
	"fmt"
)

func main () {
	conn, err := net.Dial("tcp", "52.49.91.111:2019")
	if err != nil {
		log.Fatal(err)
	}

	/**
	with previous_events as (
select a. * , (select date_time from activity oc where oc.action = 'close' and oc.user_id = a.user_id and oc.date_time < a.date_time order by oc.date_time desc limit 1) previous_close,
        (select date_time from activity oc where oc.action = 'open'  and oc.user_id = a.user_id and oc.date_time < a.date_time order by oc.date_time desc limit 1 ) previous_open
from activity a order by user_id, date_time),
session_group as (
    select
    case when
   (previous_open
   < previous_close or previous_open is null) and previous_events.action != 'open'
    then
    previous_close
    when previous_events.action = 'open'
    then
   (select date_time from activity a2 where a2.user_id = previous_events.user_id and a2.date_time
   < previous_events.date_time order by a2.date_time desc limit 1)
    else
   (select date_time from activity a2 where a2.user_id = previous_events.user_id and a2.date_time
   < previous_events.previous_open order by a2.date_time desc limit 1)
    end as session_id
   , previous_events.* from previous_events order by user_id
   , date_time
    )
select cast(user_id as integer) as user_id, cast(min(date_time) as char) as session_from, cast(max(date_time) as char) as session_to, cast(timestampdiff(second, min(date_time), max(date_time)) as integer) seconds , cast(count(*) as integer) num_actions from session_group group by user_id, session_id

;

	 */
	query := "with previous_events as ( select a. * , (select date_time from activity oc where oc.action = 'close' and oc.user_id = a.user_id and oc.date_time < a.date_time order by oc.date_time desc limit 1) previous_close,         (select date_time from activity oc where oc.action = 'open'  and oc.user_id = a.user_id and oc.date_time < a.date_time order by oc.date_time desc limit 1 ) previous_open from activity a order by user_id, date_time), session_group as (     select     case when    (previous_open    < previous_close or previous_open is null) and previous_events.action != 'open'     then     previous_close     when previous_events.action = 'open'     then    (select date_time from activity a2 where a2.user_id = previous_events.user_id and a2.date_time    < previous_events.date_time order by a2.date_time desc limit 1)     else    (select date_time from activity a2 where a2.user_id = previous_events.user_id and a2.date_time    < previous_events.previous_open order by a2.date_time desc limit 1)     end as session_id    , previous_events.* from previous_events order by user_id    , date_time     ) select cast(user_id as integer) as user_id, cast(min(date_time) as char) as session_from, cast(max(date_time) as char) as session_to, cast(timestampdiff(second, min(date_time), max(date_time)) as integer) seconds , cast(count(*) as integer) num_actions from session_group group by user_id, session_id;"
	fmt.Fprint(conn, query)
	reader := bufio.NewReader(conn)
	for true {
		status, err := reader.ReadString('\n')
		log.Println(status)
		handleError(err)

	}
}


func handleError (err error){
	if err != nil {
		log.Fatal(err)
	}
}
