import * as Knex from 'knex'

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('vc_schema_references', function (table) {
    table.uuid('schema_id').references('id').inTable('vc_schemas')
    table.string('name').notNullable()
    table.jsonb('reference_body').notNullable()
    table.dateTime('created_at').notNullable()
    table.unique(['schema_id', 'name']);
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('vc_schema_references')
}
