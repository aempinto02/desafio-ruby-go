require 'rails_helper'

RSpec.describe Services::Api::V1::Products::Details do
  let(:valid_params) { { id: 1 } }
  let(:invalid_params) { { id: 0 } }
  let(:request) { double('Request') }

  describe '#execute' do
    context 'when product with given id exists' do
      it 'returns the product' do
        product = create(:product, id: 1, price: 99.99)
        service = described_class.new(valid_params, request)

        expect(service.execute).to eq(product)
      end
    end
  end
end
