import * as Knex from 'knex'

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable('vc_schemas', function (table) {
    table.uuid('id').primary()
    table.jsonb('schema_body').notNullable()
    table.string('schema_name', 255).notNullable()
    table.string('schema_type', 255).notNullable().unique()
    table.string('version', 255).notNullable()
    table.uuid('created_by').references('id').inTable('tokens')
    table.dateTime('created_at').notNullable()
    table.dateTime('updated_at').notNullable()
    table.dateTime('deleted_at')
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('vc_schemas')
}
