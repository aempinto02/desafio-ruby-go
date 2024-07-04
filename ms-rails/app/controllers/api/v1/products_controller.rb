class Api::V1::ProductsController < ApplicationController
  def index
    default_render json: ::Services::Api::V1::Products::ListAll.new(product_params, request).execute, status: 200
  end

  def show
    default_render json: ::Services::Api::V1::Products::Details.new(product_params, request).execute, status: 200
  end

  def create
    product = upsert(product_params.merge({is_api: true}))
    
    default_render json: product, status: 201
  end

  def update
    product = upsert(product_params.merge({is_api: true}))
    
    default_render json: product, status: 201
  end

  private
  def product_params
    params&.permit(
      [
        :id,
        :name,
        :brand,
        :price,
        :description,
        :stock
      ]
    )
  end

  private
  def upsert(params)
    upsert_service = ::Services::Api::V1::Products::Upsert.new(params, request)
    product = upsert_service.execute
    
    # Send product to Kafka
    if product.persisted? # Ensure the product was saved successfully
      Karafka.producer.produce_sync(topic: 'rails-to-go', payload: product.to_json)
    end
    
    product
  end
end
