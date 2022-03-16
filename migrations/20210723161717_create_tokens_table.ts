import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable('tokens', function (table) {
        table.uuid('id').primary()
        table.string('name', 255).notNullable()
        table.string('token', 255).notNullable()
        table.string('role',255).notNullable()
        table.dateTime('created_at').notNullable()
        table.dateTime('deleted_at').nullable()
    })
}

export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists('tokens')
}

