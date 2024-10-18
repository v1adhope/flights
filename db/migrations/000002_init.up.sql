create table if not exists passengers (
  passenger_id uuid,
  first_name varchar(255) not null,
  last_name varchar(255) not null,
  middle_name varchar(255) not null,

  constraint pk_passengers_passenger_id primary key(passenger_id)
);
