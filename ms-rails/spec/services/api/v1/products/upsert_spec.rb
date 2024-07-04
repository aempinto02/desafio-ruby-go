# spec/services/api/v1/products/upsert_spec.rb
require 'rails_helper'

RSpec.describe Services::Api::V1::Products::Upsert do
  let(:valid_params) do
    {
      name: 'Valid Product',
      price: 99.99,
      brand: 'Valid Brand',
      description: 'Valid Description',
      stock: 10
    }
  end
  let(:request) { double('Request') }
  let(:upsert) { described_class.new(valid_params, request) }

  describe '#execute' do
    context 'when params[:id] > 0' do
      let!(:existing_product) { create(:product, id: 1, price: 99.99) }

      it 'updates the existing product' do
        upsert.execute
        existing_product.reload

        expect(existing_product.name).to eq('Sample Product')
        expect(existing_product.price).to eq(99.99)
      end
    end

    context 'when params[:id] is not provided' do
      it 'creates a new product' do
        expect {
          upsert.execute
        }.to change(Product, :count).by(1)
      end
    end
  end

  describe '#update_product' do
    let!(:existing_product) { create(:product, id: 1, price: 99.99) }

    it 'updates attributes of the existing product' do
      upsert.send(:update_product)
      existing_product.reload

      expect(existing_product.name).to eq('Sample Product')
      expect(existing_product.price).to eq(99.99)
    end
  end

  describe '#insert_product' do
    it 'creates a new product' do
      expect {
        upsert.send(:insert_product)
      }.to change(Product, :count).by(1)
    end
  end
end
