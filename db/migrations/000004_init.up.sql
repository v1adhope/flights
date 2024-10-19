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

create table if not exists documents (
  document_id uuid,
  type varchar(255) not null,
  number varchar(255) not null,
  passenger_id uuid not null,

  constraint pk_documents_document_id primary key(document_id),
  constraint uq_documents_type_number unique(type, number),
  constraint fk_document_passenger_passenger_id foreign key(passenger_id) references passengers(passenger_id)
);

create table if not exists passenger_ticket (
  passenger_id uuid,
  ticket_id uuid,

  constraint pk_ticket_passenger_ticket_id_passenger_id primary key (ticket_id, passenger_id),
  constraint fk_ticket_passenger_passenger_passenger_id foreign key(passenger_id) references passengers(passenger_id)  on delete cascade,
  constraint fk_ticket_passenger_tickets_ticket_id foreign key(ticket_id) references tickets(ticket_id)
)
