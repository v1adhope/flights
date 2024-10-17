create table if not exists tickets (
  ticket_id uuid not null,
  provider varchar(255) not null,
  fly_from varchar(255) not null,
  fly_to varchar(255) not null,
  fly_at timestamp not null,
  arrive_at timestamp not null,
  created_at timestamp not null,

  constraint pk_tickets_ticket_id primary key(ticket_id)
);
