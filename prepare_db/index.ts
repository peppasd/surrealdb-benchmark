import { parseArgs } from "util";
import { faker } from "@faker-js/faker";

class Table {
  constructor(public name: string) {}

  createCommand(): string {
    return `DEFINE TABLE ${this.name} SCHEMALESS PERMISSIONS NONE;\n`;
  }
}

class Customer {
  constructor(
    public id: number,
    public first_name: string,
    public last_name: string,
    public email: string,
    public country: string,
    public last_login: Date
  ) {}

  createCommand(): string {
    return `CREATE customer:${this.id} SET first_name = '${
      this.first_name
    }', last_name = '${this.last_name}', email = '${this.email}', country = '${
      this.country
    }', last_login = "${this.last_login.toISOString()}" RETURN NONE;\n`;
  }
}

class Book {
  constructor(
    public id: number,
    public title: string,
    public description: string,
    public price: number,
    public isbn: string
  ) {}

  createCommand(): string {
    return `CREATE book:${this.id} SET title = '${this.title}', description = '${this.description}', price = ${this.price}, isbn = "${this.isbn}" RETURN NONE;\n`;
  }
}

class Order {
  constructor(
    public id: number,
    public book_ids: number[],
    public customer_id: number,
    public created_at: Date,
    public processed: boolean
  ) {}

  createCommand(): string {
    return `CREATE order:${
      this.id
    } SET created_at = "${this.created_at.toISOString()}", processed = ${
      this.processed
    }, books = [${this.book_ids
      .map((bid) => `book:${bid}`)
      .join(",")}] RETURN NONE;\nRELATE customer:${
      this.customer_id
    }->ordered->order:${this.id} RETURN NONE;\n`;
  }
}

function randomInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

function parseCommandArguments(): string {
  const { values, positionals } = parseArgs({
    args: Bun.argv,
    options: {
      output: {
        type: "string",
      },
    },
    strict: true,
    allowPositionals: true,
  });

  if (values.output === undefined) {
    throw new Error("Output path is required");
  }

  return values.output;
}

async function main() {
  var outputPath = parseCommandArguments();
  const file = Bun.file(outputPath);
  const writer = file.writer();

  writer.write("OPTION IMPORT;\n");
  writer.write(new Table("customer").createCommand());
  writer.write(new Table("book").createCommand());
  writer.write(new Table("order").createCommand());
  writer.write("BEGIN TRANSACTION;\n");

  const customers_count = 200000;
  const books_count = 200000;
  const orders_count = 600000;

  for (let i = 0; i < customers_count; i++) {
    writer.write(
      new Customer(
        i,
        faker.person.firstName().replaceAll("'", " "),
        faker.person.lastName().replaceAll("'", " "),
        faker.internet.email(),
        faker.location.country().replaceAll("'", " "),
        faker.date.recent({
          days: 30,
        })
      ).createCommand()
    );
  }

  console.log("Customers created");

  for (let i = 0; i < books_count; i++) {
    writer.write(
      new Book(
        i,
        faker.commerce.productName().replaceAll("'", ""),
        faker.commerce.productDescription().replaceAll("'", ""),
        parseFloat(
          faker.commerce.price({
            max: 60,
          })
        ),
        faker.commerce.isbn()
      ).createCommand()
    );
  }

  console.log("Books created");

  for (let i = 0; i < orders_count; i++) {
    const order_amount = randomInt(1, 3);
    const book_ids = [];
    for (let j = 0; j < order_amount; j++) {
      book_ids.push(randomInt(0, books_count - 1));
    }
    writer.write(
      new Order(
        i,
        book_ids,
        randomInt(0, customers_count - 1),
        faker.date.past(),
        faker.datatype.boolean({ probability: 0.9 })
      ).createCommand()
    );
  }

  console.log("Orders created");

  writer.write("COMMIT TRANSACTION;\n");
  writer.end();
}

main();
