require 'rails_helper'

RSpec.describe Services::Api::V1::Products::ListAll do
  let(:params) { {} }
  let(:request) { double('Request') }

  describe '#execute' do
    context 'when there are products in the database' do
      it 'returns all products' do
        product1 = create(:product)
        product2 = create(:product)

        # Debugging output
        puts "Product 1: #{product1.attributes}"
        puts "Product 2: #{product2.attributes}"

        service = described_class.new(params, request)
        result = service.execute

        expect(result).to include(product1, product2)
      end
    end

    context 'when there are no products in the database' do
      it 'returns an empty ActiveRecord::Relation' do
        service = described_class.new(params, request)
        result = service.execute

        expect(result).to be_empty
      end
    end
  end
end
