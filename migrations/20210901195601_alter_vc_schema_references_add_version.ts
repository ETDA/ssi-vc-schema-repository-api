import * as Knex from 'knex'

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('vc_schema_references', function (table) {
	  table.dropForeign(["schema_id"])
	  table.dropUnique(["schema_id", "name"])
	  table.string("version", 255).notNullable().defaultTo("")
	  table.unique(["schema_id", "name", "version"])
	  table.uuid('schema_id').references('id').inTable('vc_schemas').alter()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('vc_schema_references', function (table) {
	  table.dropForeign(["schema_id"])
    table.dropUnique(["schema_id", "name", "version"])
    table.dropColumn("version")
    table.unique(['schema_id', 'name'])
    table.uuid('schema_id').references('id').inTable('vc_schemas').alter()
  })
}
