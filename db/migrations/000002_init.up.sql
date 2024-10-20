create table if not exists tickets (
  ticket_id uuid,
  provider varchar(255) not null,
  fly_from varchar(255) not null,
  fly_to varchar(255) not null,
  fly_at timestamp with time zone not null,
  arrive_at timestamp with time zone not null,
  created_at timestamp with time zone not null,

  constraint pk_tickets_ticket_id primary key(ticket_id)
);

create table if not exists passengers (
  passenger_id uuid,
  first_name varchar(255) not null,
  last_name varchar(255) not null,
  middle_name varchar(255) not null,

  constraint pk_passengers_passenger_id primary key(passenger_id)
);
