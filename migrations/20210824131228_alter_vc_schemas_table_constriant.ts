import * as Knex from 'knex'

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('vc_schemas', function (table) {
    table.string('schema_name', 255).notNullable().unique().alter()
    table.dropUnique(['schema_type'])
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('vc_schemas', function (table) {
	  table.string('schema_type', 255).notNullable().unique().alter()
	  table.dropUnique(['schema_name'])
  })
}
