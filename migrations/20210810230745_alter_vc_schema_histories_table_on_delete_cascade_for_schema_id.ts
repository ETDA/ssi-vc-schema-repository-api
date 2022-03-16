import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("vc_schema_histories", function(table) {
    table.uuid('schema_id').notNullable().references('id').inTable('vc_schemas').onDelete('CASCADE').alter()
  })
}


export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("vc_schema_histories", function(table) {
    table.dropForeign(["schema_id"])
  })
}

