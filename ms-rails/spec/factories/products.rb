FactoryBot.define do
  factory :product do
    name { "Sample Product" }
    brand { "Sample Brand" }
    price { 100.11 }
    description { "Sample Description" }
    stock { 10 }
    created_at { "2024-07-02T20:11:01.077Z" }
    updated_at { "2024-07-02T20:11:01.077Z" }
  end
end
